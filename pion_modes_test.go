// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0

package opus

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"slices"
	"testing"
	"unsafe"
)

func TestPionPLCContinuesAfterRedundancy(t *testing.T) {
	// Three SILK packets end with CELT redundancy, followed by three lost
	// 2.5 ms packets and two received CELT packets. All short losses must remain
	// in CELT state, rather than falling back to SILK after the first loss.
	var fixture struct {
		Rate, Channels int
		Steps          []struct {
			Packet  string
			Samples int
		}
	}
	data, err := os.ReadFile("testdata/pion-redundancy.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	d, err := NewDecoderWithOutput(fixture.Rate, fixture.Channels)
	if err != nil {
		t.Fatal(err)
	}
	for i, step := range fixture.Steps {
		out := make([]int16, step.Samples*fixture.Channels)
		if step.Packet == "" {
			if err := d.DecodePLC(out); err != nil {
				t.Fatalf("PLC step %d: %v", i, err)
			}
		} else {
			packet, err := hex.DecodeString(step.Packet)
			if err != nil {
				t.Fatal(err)
			}
			n, err := d.DecodeToInt16(packet, out)
			if err != nil || n != step.Samples {
				t.Fatalf("decode step %d: %d %v", i, n, err)
			}
		}
	}
}

// Compare the slice API with direct codec calls, including Hybrid packets that
// the Pion-compatible CELT/SILK encoder entry points do not produce.
func TestPionDecoderModes(t *testing.T) {
	tls := NewTLS()
	defer tls.Close()
	errPtr, va, input, packetPtr, output := Malloc(4), Malloc(8), Malloc(1920*2), Malloc(1500), Malloc(1920*4)
	defer func() {
		for _, p := range []uintptr{errPtr, va, input, packetPtr, output} {
			Free(p)
		}
	}()
	for _, mode := range []int32{MODE_SILK_ONLY, MODE_HYBRID, MODE_CELT_ONLY} {
		for _, channels := range []int32{1, 2} {
			enc := opus_encoder_create(tls, 48000, channels, OPUS_APPLICATION_AUDIO, errPtr)
			if enc == 0 {
				t.Fatal("encoder create")
			}
			if r := opus_encoder_ctl(tls, enc, OPUS_SET_FORCE_MODE_REQUEST, VaList(va, mode)); r != 0 {
				t.Fatal(r)
			}
			if r := opus_encoder_ctl(tls, enc, OPUS_SET_BITRATE_REQUEST, VaList(va, int32(64000))); r != 0 {
				t.Fatal(r)
			}
			if r := opus_encoder_ctl(tls, enc, OPUS_SET_FORCE_CHANNELS_REQUEST, VaList(va, channels)); r != 0 {
				t.Fatal(r)
			}
			if mode == MODE_SILK_ONLY {
				if r := opus_encoder_ctl(tls, enc, OPUS_SET_BANDWIDTH_REQUEST, VaList(va, int32(OPUS_BANDWIDTH_WIDEBAND))); r != 0 {
					t.Fatal(r)
				}
			}
			pcm := unsafe.Slice((*int16)(unsafe.Pointer(input)), int(960*channels))
			for i := range pcm {
				pcm[i] = int16(10000 * math.Sin(float64(i/int(channels))*.08))
			}
			n := opus_encode(tls, enc, input, 960, packetPtr, 1500)
			opus_encoder_destroy(tls, enc)
			if n <= 0 {
				t.Fatalf("encode mode %d: %d", mode, n)
			}
			packet := slices.Clone(unsafe.Slice((*byte)(unsafe.Pointer(packetPtr)), int(n)))
			if got := opus_packet_get_mode(tls, packetPtr); got != mode {
				t.Fatalf("fixture mode: %d, want %d", got, mode)
			}
			for _, rate := range []int32{8000, 12000, 16000, 24000, 48000} {
				for _, outChannels := range []int32{1, 2} {
					d, err := NewDecoderWithOutput(int(rate), int(outChannels))
					if err != nil {
						t.Fatal(err)
					}
					ref := opus_decoder_create(tls, rate, outChannels, errPtr)
					if ref == 0 {
						t.Fatal("decoder create")
					}
					got := make([]float32, int(rate/50*outChannels))
					for range 3 {
						count, err := d.DecodeToFloat32(packet, got)
						wantCount := opus_decode_float(tls, ref, packetPtr, n, output, rate/50, 0)
						if err != nil || count != int(wantCount) || count != int(rate/50) {
							t.Fatalf("mode=%d rate=%d channels=%d: %d %d %v", mode, rate, outChannels, count, wantCount, err)
						}
						want := unsafe.Slice((*float32)(unsafe.Pointer(output)), count*int(outChannels))
						if !slices.Equal(got, want) {
							t.Fatalf("PCM differs: mode=%d rate=%d channels=%d", mode, rate, outChannels)
						}
					}
					opus_decoder_destroy(tls, ref)
				}
			}
		}
	}
}

func TestPionMalformedPackets(t *testing.T) {
	packets := [][]byte{
		nil,
		{0xf9, 0, 0, 0}, // code 1, odd payload size
		{0xfa},          // code 2, missing length
		{0xfb, 0},       // code 3, zero frames
		{0xfb, 7},       // seven 20 ms frames exceed 120 ms
		{0xfb, 0x41},    // padding flag, missing padding length
		append([]byte{0xf8}, make([]byte, 1276)...), // oversized frame
	}
	d := NewDecoder()
	for _, packet := range packets {
		out := make([]int16, 5760)
		for i := range out {
			out[i] = 2345
		}
		if _, err := d.DecodeToInt16(packet, out); err == nil {
			t.Fatalf("malformed packet accepted: %x", packet)
		}
		for _, sample := range out {
			if sample != 2345 {
				t.Fatal("failed decode wrote output")
			}
		}
	}
}

func TestPionShortCELTPLC(t *testing.T) {
	e, err := NewEncoder()
	if err != nil {
		t.Fatal(err)
	}
	packet := make([]byte, 1500)
	n, err := e.EncodeFloat32(make([]float32, 960), packet)
	if err != nil {
		t.Fatal(err)
	}
	d := NewDecoder()
	if _, err := d.DecodeToInt16(packet[:n], make([]int16, 960)); err != nil {
		t.Fatal(err)
	}
	for _, count := range []int{120, 240, 480, 960} {
		if err := d.DecodePLC(make([]int16, count)); err != nil {
			t.Fatalf("%d samples: %v", count, err)
		}
	}
	if err := d.DecodePLC(make([]int16, 1920)); err == nil {
		t.Fatal("long CELT PLC accepted")
	}
}
