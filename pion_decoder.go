// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0

package opus

import (
	"encoding/binary"
	"math"
	"runtime"
	"unsafe"
)

// Decoder provides the pion/opus PCM decoding API. PCM is interleaved when
// channels is two. A decoder must not be used concurrently or copied after use.
// Native resources are reclaimed automatically when the decoder becomes unreachable.
type Decoder struct {
	native               *nativeCodec
	sampleRate, channels int
	softClipMem          softClipMemory
}

// NewDecoder creates a 48 kHz mono decoder.
func NewDecoder() Decoder {
	d, err := NewDecoderWithOutput(48000, 1)
	if err != nil {
		panic(err)
	}
	return d
}

// NewDecoderWithOutput creates a decoder at 8, 12, 16, 24 or 48 kHz with one or two channels.
func NewDecoderWithOutput(sampleRate, channels int) (Decoder, error) {
	var d Decoder
	if err := d.Init(sampleRate, channels); err != nil {
		return Decoder{}, err
	}
	return d, nil
}

// Init initializes or resets a decoder. Invalid settings leave its state unchanged.
func (d *Decoder) Init(sampleRate, channels int) error {
	if err := validOutput(sampleRate, channels); err != nil {
		return err
	}
	n, err := newNativeCodec(sampleRate, channels, false)
	if err != nil {
		return err
	}
	*d = Decoder{native: n, sampleRate: sampleRate, channels: channels}
	return nil
}

// prepare validates a packet and checks output capacity before advancing codec state.
func (d *Decoder) prepare(in []byte, outputSamples int) (int, error) {
	if err := validOutput(d.sampleRate, d.channels); err != nil {
		return 0, err
	}
	if len(in) == 0 {
		return 0, errTooShortForTableOfContentsHeader
	}
	if uint64(len(in)) > math.MaxInt32 {
		return 0, errMalformedPacket
	}
	defer runtime.KeepAlive(d.native)
	r := d.native.r
	if err := r.buffers(len(in), d.sampleRate*120/1000*d.channels*4); err != nil {
		return 0, err
	}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(r.input)), len(in)), in)
	count := opus_packet_parse(r.tls, r.input, int32(len(in)), 0, 0, r.scratch, 0)
	if count <= 0 {
		return 0, errMalformedPacket
	}
	for _, size := range unsafe.Slice((*int16)(unsafe.Pointer(r.scratch)), int(count)) {
		if size < 0 || size > 1275 {
			return 0, errMalformedPacket
		}
	}
	n := opus_decoder_get_nb_samples(r.tls, r.state, r.input, int32(len(in)))
	if n <= 0 {
		return 0, errMalformedPacket
	}
	if outputSamples < int(n)*d.channels {
		return 0, errOutBufferTooSmall
	}
	return int(n), nil
}

func (d *Decoder) decodePCM(in []byte, outputSamples int, integer bool) ([]float32, int, error) {
	n, err := d.prepare(in, outputSamples)
	if err != nil {
		return nil, 0, err
	}
	defer runtime.KeepAlive(d.native)
	r := d.native.r
	count := opus_decode_float(r.tls, r.state, r.input, int32(len(in)), r.output, int32(n), 0)
	if err := codecError(count); err != nil {
		return nil, 0, err
	}
	out := unsafe.Slice((*float32)(unsafe.Pointer(r.output)), int(count)*d.channels)
	if integer {
		softClip(out, d.channels, &d.softClipMem)
	} else {
		d.softClipMem = softClipMemory{}
	}
	return out, int(count), nil
}

func pionInt16(sample float32) int16 {
	v := math.RoundToEven(float64(sample * 32768))
	return int16(min(32767, max(-32768, v)))
}

// Decode writes S16LE PCM and returns the packet bandwidth and encoded stereo flag.
func (d *Decoder) Decode(in, out []byte) (Bandwidth, bool, error) {
	defer runtime.KeepAlive(d.native)
	pcm, _, err := d.decodePCM(in, len(out)/2, true)
	if err != nil {
		return 0, false, err
	}
	for i, v := range pcm {
		binary.LittleEndian.PutUint16(out[2*i:], uint16(pionInt16(v)))
	}
	return packetBandwidth(in[0]), in[0]&4 != 0, nil
}

// DecodeFloat32 writes float PCM and returns the packet bandwidth and encoded stereo flag.
func (d *Decoder) DecodeFloat32(in []byte, out []float32) (Bandwidth, bool, error) {
	if _, err := d.DecodeToFloat32(in, out); err != nil {
		return 0, false, err
	}
	return packetBandwidth(in[0]), in[0]&4 != 0, nil
}

// DecodeToInt16 writes signed PCM and returns the number of samples per channel.
func (d *Decoder) DecodeToInt16(in []byte, out []int16) (int, error) {
	defer runtime.KeepAlive(d.native)
	pcm, count, err := d.decodePCM(in, len(out), true)
	if err != nil {
		return 0, err
	}
	for i, v := range pcm {
		out[i] = pionInt16(v)
	}
	return count, nil
}

// DecodeToFloat32 writes float PCM and returns the number of samples per channel.
func (d *Decoder) DecodeToFloat32(in []byte, out []float32) (int, error) {
	defer runtime.KeepAlive(d.native)
	pcm, count, err := d.decodePCM(in, len(out), false)
	if err != nil {
		return 0, err
	}
	copy(out, pcm)
	return count, nil
}

// DecodePLC conceals a missing packet of the duration specified by out. Before
// the first packet only 20 ms is accepted. SILK also accepts 40/60 ms, CELT
// accepts 2.5/5/10 ms; unsupported lengths leave state and output unchanged.
func (d *Decoder) DecodePLC(out []int16) error {
	if err := validOutput(d.sampleRate, d.channels); err != nil {
		return err
	}
	defer runtime.KeepAlive(d.native)
	r := d.native.r
	state := (*OpusDecoder)(unsafe.Pointer(r.state))
	mode := state.Fprev_mode
	if mode != 0 && state.Fprev_redundancy != 0 {
		mode = MODE_CELT_ONLY
	}
	base := d.sampleRate / 50 * d.channels
	valid := len(out) == base || mode == MODE_CELT_ONLY && (len(out) == base/8 || len(out) == base/4 || len(out) == base/2) || mode == MODE_SILK_ONLY && (len(out) == 2*base || len(out) == 3*base)
	if !valid {
		return errInvalidPLCFrameSize
	}
	if mode == 0 {
		clear(out)
		return nil
	}
	if err := r.buffers(0, len(out)*4); err != nil {
		return err
	}
	if state.Fprev_redundancy != 0 {
		// Pion conceals from the final CELT redundancy while retaining the
		// received packet's mode for the next real packet's transition logic.
		previousMode, previousRedundancy := state.Fprev_mode, state.Fprev_redundancy
		state.Fprev_mode, state.Fprev_redundancy = MODE_CELT_ONLY, 0
		defer func() { state.Fprev_mode, state.Fprev_redundancy = previousMode, previousRedundancy }()
	}
	// The RFC 6716 core caps null-packet calls to the previous TOC duration.
	// SILK's public PLC contract allows a whole lost 40/60 ms packet instead.
	previousFrame := state.Fframe_size
	if mode == MODE_SILK_ONLY {
		state.Fframe_size = int32(len(out) / d.channels)
		defer func() { state.Fframe_size = previousFrame }()
	}
	for offset := 0; offset < len(out); {
		n := opus_decode_float(r.tls, r.state, 0, 0, r.output+uintptr(offset)*4, int32((len(out)-offset)/d.channels), 0)
		if err := codecError(n); err != nil {
			return err
		}
		if n == 0 {
			return errInvalidPLCFrameSize
		}
		offset += int(n) * d.channels
	}
	for i, v := range unsafe.Slice((*float32)(unsafe.Pointer(r.output)), len(out)) {
		out[i] = pionInt16(v)
	}
	return nil
}
