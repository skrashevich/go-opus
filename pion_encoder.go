// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0

package opus

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"slices"
	"unsafe"
)

// Encoder provides the pion/opus PCM encoding API. Encode and EncodeFloat32
// accept one 20 ms frame at 48 kHz. EncodeSILK accepts mono 20/40/60 ms frames
// at the selected SILK bandwidth's sample rate. Do not use an encoder concurrently.
type Encoder struct {
	native                                              *nativeCodec
	silk                                                [3]*nativeCodec
	sampleRate, channels, bitrate, complexity, lossRate int
	application                                         Application
	vbr, constrainedVBR                                 bool
	bandwidth, maxBandwidth                             Bandwidth
}

// EncoderOption configures an encoder; options are applied in order.
type EncoderOption func(*Encoder) error

// WithSampleRate sets the PCM input rate. The compatible encoder API uses 48 kHz.
func WithSampleRate(rate int) EncoderOption {
	return func(e *Encoder) error {
		if rate != 48000 {
			return errInvalidSampleRate
		}
		e.sampleRate = rate
		return nil
	}
}

// WithChannels sets the PCM channel count to one or two.
func WithChannels(channels int) EncoderOption {
	return func(e *Encoder) error {
		if channels != 1 && channels != 2 {
			return errInvalidChannelCount
		}
		e.channels = channels
		return nil
	}
}

// WithBitrate sets a target of 6000 through 510000 bits per second.
func WithBitrate(bps int) EncoderOption { return func(e *Encoder) error { return e.SetBitrate(bps) } }

// WithComplexity sets encoder complexity, from zero through ten.
func WithComplexity(complexity int) EncoderOption {
	return func(e *Encoder) error { return e.SetComplexity(complexity) }
}

// WithApplication records the tuning profile, as in pion/opus.
func WithApplication(app Application) EncoderOption {
	return func(e *Encoder) error { return e.SetApplication(app) }
}

// WithVBR selects variable bitrate encoding.
func WithVBR(vbr bool) EncoderOption { return func(e *Encoder) error { e.SetVBR(vbr); return nil } }

// WithConstrainedVBR selects constrained variable bitrate encoding.
func WithConstrainedVBR(cvbr bool) EncoderOption {
	return func(e *Encoder) error { e.SetConstrainedVBR(cvbr); return nil }
}

// WithBandwidth fixes the CELT bandwidth. Auto and Mediumband are rejected, matching pion/opus.
func WithBandwidth(bw Bandwidth) EncoderOption {
	return func(e *Encoder) error { return e.SetBandwidth(bw) }
}

// WithMaxBandwidth caps automatic CELT bandwidth selection.
func WithMaxBandwidth(bw Bandwidth) EncoderOption {
	return func(e *Encoder) error { return e.SetMaxBandwidth(bw) }
}

// NewEncoder creates a 48 kHz mono encoder at 24 kbit/s, complexity five,
// ApplicationAudio, CBR, and constrained VBR enabled, matching pion/opus defaults.
func NewEncoder(opts ...EncoderOption) (*Encoder, error) {
	e := &Encoder{sampleRate: 48000, channels: 1, bitrate: 24000, complexity: 5,
		application: ApplicationAudio, constrainedVBR: true, maxBandwidth: BandwidthFullband}
	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	n, err := newNativeCodec(e.sampleRate, e.channels, true)
	if err != nil {
		return nil, err
	}
	e.native = n
	if err := e.configure(n, MODE_CELT_ONLY, e.channels); err != nil {
		return nil, err
	}
	return e, nil
}

func boolInt(v bool) int32 {
	if v {
		return 1
	}
	return 0
}

func (e *Encoder) configure(n *nativeCodec, mode, channels int) error {
	for _, setting := range [][2]int32{
		{OPUS_SET_BITRATE_REQUEST, int32(e.bitrate)},
		{OPUS_SET_COMPLEXITY_REQUEST, int32(e.complexity)},
		{OPUS_SET_VBR_REQUEST, boolInt(e.vbr)},
		{OPUS_SET_VBR_CONSTRAINT_REQUEST, boolInt(e.constrainedVBR)},
		{OPUS_SET_PACKET_LOSS_PERC_REQUEST, int32(e.lossRate)},
		{OPUS_SET_FORCE_MODE_REQUEST, int32(mode)},
		{OPUS_SET_FORCE_CHANNELS_REQUEST, int32(channels)},
	} {
		if err := n.ctl(setting[0], setting[1]); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) setCtl(request, value int32) error {
	if e.native != nil {
		return e.native.ctl(request, value)
	}
	return nil
}

// SetBitrate updates the bitrate without resetting encoder state.
func (e *Encoder) SetBitrate(bps int) error {
	if bps < 6000 || bps > 510000 {
		return fmt.Errorf("%w: %d", errBitrateOutOfRange, bps)
	}
	if err := e.setCtl(OPUS_SET_BITRATE_REQUEST, int32(bps)); err != nil {
		return err
	}
	e.bitrate = bps
	return nil
}

// SetComplexity updates encoder complexity, from zero through ten.
func (e *Encoder) SetComplexity(complexity int) error {
	if complexity < 0 || complexity > 10 {
		return fmt.Errorf("%w: %d", errInvalidComplexity, complexity)
	}
	if err := e.setCtl(OPUS_SET_COMPLEXITY_REQUEST, int32(complexity)); err != nil {
		return err
	}
	e.complexity = complexity
	return nil
}

// SetApplication records the application profile. Like pion/opus, it does not change other controls.
func (e *Encoder) SetApplication(app Application) error {
	switch app {
	case ApplicationAudio, ApplicationVoIP, ApplicationRestrictedLowDelay:
	default:
		return fmt.Errorf("%w: %d", errInvalidApplication, app)
	}
	e.application = app
	return nil
}

// SetVBR enables or disables variable bitrate encoding.
func (e *Encoder) SetVBR(vbr bool) {
	// Both boolean values are accepted by the native control.
	_ = e.setCtl(OPUS_SET_VBR_REQUEST, boolInt(vbr))
	e.vbr = vbr
}

// SetConstrainedVBR enables or disables constrained VBR.
func (e *Encoder) SetConstrainedVBR(cvbr bool) {
	_ = e.setCtl(OPUS_SET_VBR_CONSTRAINT_REQUEST, boolInt(cvbr))
	e.constrainedVBR = cvbr
}

// SetLossRate sets expected packet loss, from zero through 100 percent.
func (e *Encoder) SetLossRate(rate int) error {
	if rate < 0 || rate > 100 {
		return fmt.Errorf("%w: %d", errInvalidLossRate, rate)
	}
	if err := e.setCtl(OPUS_SET_PACKET_LOSS_PERC_REQUEST, int32(rate)); err != nil {
		return err
	}
	e.lossRate = rate
	return nil
}

func validCELTBandwidth(bw Bandwidth) error {
	if bw < BandwidthNarrowband || bw > BandwidthFullband || bw == BandwidthMediumband {
		return fmt.Errorf("%w: %d", errInvalidBandwidth, bw)
	}
	return nil
}

// SetBandwidth fixes the CELT bandwidth, overriding automatic selection.
func (e *Encoder) SetBandwidth(bw Bandwidth) error {
	if err := validCELTBandwidth(bw); err != nil {
		return err
	}
	e.bandwidth = bw
	return nil
}

// SetMaxBandwidth caps automatic bandwidth selection.
func (e *Encoder) SetMaxBandwidth(bw Bandwidth) error {
	if err := validCELTBandwidth(bw); err != nil {
		return err
	}
	e.maxBandwidth = bw
	return nil
}

func (e *Encoder) Application() Application { return e.application }
func (e *Encoder) VBR() bool                { return e.vbr }
func (e *Encoder) Complexity() int          { return e.complexity }
func (e *Encoder) ConstrainedVBR() bool     { return e.constrainedVBR }
func (e *Encoder) LossRate() int            { return e.lossRate }
func (e *Encoder) Bandwidth() Bandwidth     { return e.bandwidth }
func (e *Encoder) MaxBandwidth() Bandwidth  { return e.maxBandwidth }

func (e *Encoder) selectedBandwidth() Bandwidth {
	if e.bandwidth != BandwidthAuto {
		return e.bandwidth
	}
	rate := e.bitrate
	if !e.vbr {
		rate -= rate / 12
	}
	rate = rate * (90 + e.complexity) / 100
	bw := BandwidthFullband
	switch {
	case rate < 9000:
		bw = BandwidthNarrowband
	case rate < 13500:
		bw = BandwidthWideband
	case rate < 14000:
		bw = BandwidthSuperwideband
	}
	return min(bw, e.maxBandwidth)
}

func (e *Encoder) prepareFrame(samples, outputBytes int) error {
	if e.sampleRate != 48000 || e.native == nil {
		return errInvalidSampleRate
	}
	if samples != 960*e.channels {
		return errInvalidFrameSize
	}
	if outputBytes < e.bitrate/400 {
		return errOutBufferTooSmall
	}
	// An explicit bandwidth takes precedence over the automatic cap.
	bw := int32(e.selectedBandwidth()) + 1100
	if err := e.native.ctl(OPUS_SET_MAX_BANDWIDTH_REQUEST, OPUS_BANDWIDTH_FULLBAND); err != nil {
		return err
	}
	if err := e.native.ctl(OPUS_SET_BANDWIDTH_REQUEST, bw); err != nil {
		return err
	}
	return e.native.r.buffers(samples*4, min(outputBytes, 1276))
}

// Encode encodes interleaved S16LE PCM into a single 20 ms CELT packet.
func (e *Encoder) Encode(in, out []byte) (int, error) {
	if len(in)%2 != 0 {
		return 0, errInvalidInputLength
	}
	if err := e.prepareFrame(len(in)/2, len(out)); err != nil {
		return 0, err
	}
	defer runtime.KeepAlive(e.native)
	r := e.native.r
	pcm := unsafe.Slice((*int16)(unsafe.Pointer(r.input)), len(in)/2)
	for i := range pcm {
		pcm[i] = int16(binary.LittleEndian.Uint16(in[2*i:]))
	}
	n := opus_encode(r.tls, r.state, r.input, 960, r.output, int32(min(len(out), 1276)))
	return copyEncodedPacket(r, n, out)
}

// EncodeFloat32 encodes interleaved float PCM into a single 20 ms CELT packet.
func (e *Encoder) EncodeFloat32(in []float32, out []byte) (int, error) {
	if err := e.prepareFrame(len(in), len(out)); err != nil {
		return 0, err
	}
	defer runtime.KeepAlive(e.native)
	r := e.native.r
	copy(unsafe.Slice((*float32)(unsafe.Pointer(r.input)), len(in)), in)
	n := opus_encode_float(r.tls, r.state, r.input, 960, r.output, int32(min(len(out), 1276)))
	return copyEncodedPacket(r, n, out)
}

func copyEncodedPacket(r *codecResources, n int32, out []byte) (int, error) {
	if err := codecError(n); err != nil {
		return 0, err
	}
	copy(out, unsafe.Slice((*byte)(unsafe.Pointer(r.output)), int(n)))
	return int(n), nil
}

// EncodeSILK encodes mono PCM at 8/12/16 kHz for narrow/medium/wide bandwidth.
// Frame durations are 20, 40 and 60 ms. Output capacity failures preserve state.
func (e *Encoder) EncodeSILK(pcm []int16, bandwidth Bandwidth, out []byte) (int, error) {
	if bandwidth < BandwidthNarrowband || bandwidth > BandwidthWideband {
		return 0, errInvalidBandwidth
	}
	unit := bandwidth.SampleRate() / 50
	if len(pcm) != unit && len(pcm) != 2*unit && len(pcm) != 3*unit {
		return 0, errInvalidFrameSize
	}
	idx := bandwidth - 1
	if e.silk[idx] == nil {
		n, err := newNativeCodec(bandwidth.SampleRate(), 1, true)
		if err != nil {
			return 0, err
		}
		e.silk[idx] = n
	}
	native := e.silk[idx]
	defer runtime.KeepAlive(native)
	if err := e.configure(native, MODE_SILK_ONLY, 1); err != nil {
		return 0, err
	}
	if err := native.ctl(OPUS_SET_BANDWIDTH_REQUEST, int32(bandwidth)+1100); err != nil {
		return 0, err
	}
	r := native.r
	if err := r.buffers(len(pcm)*2, 1276); err != nil {
		return 0, err
	}
	copy(unsafe.Slice((*int16)(unsafe.Pointer(r.input)), len(pcm)), pcm)
	// The encoder state contains offsets into its contiguous allocation, so a
	// byte snapshot can roll back a packet that does not fit the caller's buffer.
	state := unsafe.Slice((*byte)(unsafe.Pointer(r.state)), int(opus_encoder_get_size(r.tls, 1)))
	snapshot := slices.Clone(state)
	count := opus_encode(r.tls, r.state, r.input, int32(len(pcm)), r.output, 1276)
	if err := codecError(count); err != nil {
		copy(state, snapshot)
		return 0, err
	}
	if int(count) > len(out) {
		copy(state, snapshot)
		return 0, errOutBufferTooSmall
	}
	return copyEncodedPacket(r, count, out)
}
