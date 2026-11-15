// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

package opus

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"runtime"
	"testing"
	"unsafe"
)

// heapSink forces codec buffers onto the heap. The codec receives raw
// uintptrs, which do not keep a goroutine stack slot in place: if the stack
// grows during a call, writes through the stale address are lost.
var heapSink any

func cBuf[T any](n int) []T {
	s := make([]T, n)
	heapSink = s
	return s
}

func TestCodecOutput(t *testing.T) {
	for _, tc := range []struct {
		name        string
		rate        int32
		frameSize   int32
		bitrate     int32 // 0 keeps the encoder default
		channels    int32
		application int32
		signal      int // 0 sine, 1 noise, 2 chirp with noise bursts
		frames      int
		want        string
		wantPackets string // encoder bitstream remains RFC 6716-compatible
	}{
		{"stereo audio", 48000, 960, 0, 2, OPUS_APPLICATION_AUDIO, 0, 10, "7100148eac64468b99fa66324ed57a454acd1d4ae257530dfc1c46e185883192", "736e53bdf9eaf38098d7a2b4a969af9f1ed451c13dedc634be1b3f51795de05a"},
		{"mono voip", 48000, 960, 0, 1, OPUS_APPLICATION_VOIP, 0, 10, "ad1d8ea999f460474691cb2164cce788ef52e4433f8aca9c7fa55a28d5fb8650", "00f1f939fe3b3f79a904e2bd2a0d71db47f546869c9572b0eaf6aaa50bd9c251"},
		{"stereo noise", 48000, 960, 0, 2, OPUS_APPLICATION_AUDIO, 1, 10, "d55f39654d51200301f7e2bce6aacb712f20af0a9322e350a5f60976e38c55fd", "61166aafcba815c05c2bd011ecbb13e609373da1ed2efb33af7d3fe794110e55"},
		{"silk 16k mono voip", 16000, 320, 16000, 1, OPUS_APPLICATION_VOIP, 0, 10, "8f3fefa3be77c2fae500056bd583497a2ad3176664b1b099e79d51c22fb37672", "d8e4153afb2e71c56c7b3e5d5b5d86c3c9b1435297f60e1e88c214f6349d3e1f"},
		{"silk 8k mono noise", 8000, 160, 12000, 1, OPUS_APPLICATION_VOIP, 1, 10, "4bb1df117f6ef68400a0d759b789c4f84f1a836e9cdccb30fb261333610477e1", "bc0dcc7e6c3cab58c193613ade10899b4c400b9b25b726a97e0571b44e039f46"},
		{"hybrid stereo 24k", 48000, 960, 24000, 2, OPUS_APPLICATION_AUDIO, 0, 10, "8dffbb7ec67581fac279342ca55433f165460e13977e300c4cbda3610c267be1", "d181e524df34b07e5dbf1fe800874eb6d6baa1edbd3abacae7b1282897edb5e4"},
		{"celt stereo 10ms noise", 48000, 480, 96000, 2, OPUS_APPLICATION_AUDIO, 1, 10, "70f7cd94e110de21b2566ea866fd29067eff6be9143ba900a5635e14c03a676c", "6a51b7edb4a0c763a751ebcd496caf74ba8d86d632a27892887941e1d8eb7e54"},
		{"mix stereo audio long", 48000, 960, 0, 2, OPUS_APPLICATION_AUDIO, 2, 200, "57d607b522c5a3e7d8debd405c97bc04963c3c97427d8c93f157ec007e8f6f23", "d2950a4cfe6de1834bca27c396d9bd222b4646caca504522ce97af8ec93fd0c5"},
		{"mix stereo 10ms 64k long", 48000, 480, 64000, 2, OPUS_APPLICATION_AUDIO, 2, 200, "bcfb5ec04c459fe04e1e1e509010be385a0302efccb0ad119f7bcc2b23fc88d8", "956b2468e0d8620f36fcf065f48c82dd7fdd95020b7add1ddbba57831d2e950b"},
		{"mix hybrid mono long", 48000, 960, 20000, 1, OPUS_APPLICATION_VOIP, 2, 200, "929b9a8cf43fa5575e0b6195fddcecaa56cb7a9dc6adebe1b0ac753c7b4ad30b", "ac15115837c977f25d84986e24cc68435b302cb16e96e39afb58f2201de4130e"},
		{"mix silk 16k long", 16000, 320, 14000, 1, OPUS_APPLICATION_VOIP, 2, 200, "4ed553920c3b613b79c7699f7627e54af1d4e577cf14ebe9faa0bc340d5b9ef4", "31f0aaf060770a9de7cb7b897758861876eb5d0feaf88f39a66d981df2b2e1cb"},
		{"mix 24k stereo 40ms long", 24000, 960, 48000, 2, OPUS_APPLICATION_AUDIO, 2, 100, "322197c02cb1dd2d352d2d5a2dd1b5a80c02aaafe456e72a9169537e2372b5a0", "97029bac631f0910fee7623243940c505aa1fef91ba50a9c945e532f1ea40c7c"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tls := crtNewTLS()
			defer tls.Close()
			status := &cBuf[int32](1)[0]
			enc := opus_encoder_create(tls, tc.rate, tc.channels, tc.application, uintptr(unsafe.Pointer(status)))
			if *status != OPUS_OK || enc == 0 {
				t.Fatalf("encoder init: %d", *status)
			}
			defer opus_encoder_destroy(tls, enc)
			if tc.bitrate != 0 {
				va := tls.Alloc(16)
				rc := opus_encoder_ctl(tls, enc, OPUS_SET_BITRATE_REQUEST, crtVaList(va, tc.bitrate))
				tls.Free(16)
				if rc != OPUS_OK {
					t.Fatalf("set bitrate: %d", rc)
				}
			}
			dec := opus_decoder_create(tls, tc.rate, tc.channels, uintptr(unsafe.Pointer(status)))
			if *status != OPUS_OK || dec == 0 {
				t.Fatalf("decoder init: %d", *status)
			}
			defer opus_decoder_destroy(tls, dec)
			pcm := cBuf[int16](int(tc.frameSize * tc.channels))
			out := cBuf[int16](int(tc.frameSize * tc.channels))
			packet := cBuf[byte](1275)
			h := sha256.New()
			encoded := sha256.New()
			pcmLE := make([]byte, 0, 2*len(out))
			seed := uint32(1)
			for frame := range tc.frames {
				for i := range pcm {
					switch tc.signal {
					case 1:
						seed = 1664525*seed + 1013904223
						pcm[i] = int16(seed>>16) / 2
					case 2:
						// Rising chirp, amplitude modulated, with periodic noise bursts.
						t := float64(frame*len(pcm)+i) / float64(tc.rate) / float64(tc.channels)
						v := 8000 * math.Sin(2*math.Pi*(100+400*t)*t) * (0.6 + 0.4*math.Sin(2*math.Pi*1.3*t))
						seed = 1664525*seed + 1013904223
						if frame%17 < 3 {
							v += float64(int16(seed>>16)) / 3
						} else {
							v += float64(int16(seed>>16)) / 200
						}
						pcm[i] = int16(v)
					default:
						pcm[i] = int16(10000 * math.Sin(2*math.Pi*440*float64(frame*len(pcm)+i)/float64(tc.rate)))
					}
				}
				n := opus_encode(tls, enc, uintptr(unsafe.Pointer(&pcm[0])), tc.frameSize, uintptr(unsafe.Pointer(&packet[0])), int32(len(packet)))
				if n <= 0 {
					t.Fatalf("encode frame %d: %d", frame, n)
				}
				h.Write(packet[:n])
				encoded.Write(packet[:n])
				got := opus_decode(tls, dec, uintptr(unsafe.Pointer(&packet[0])), n, uintptr(unsafe.Pointer(&out[0])), tc.frameSize, 0)
				if got <= 0 || got > tc.frameSize {
					t.Fatalf("decode frame %d: %d", frame, got)
				}
				// Hash PCM little-endian so the digests hold on big-endian targets too.
				pcmLE = pcmLE[:0]
				for _, s := range out[:got*tc.channels] {
					pcmLE = binary.LittleEndian.AppendUint16(pcmLE, uint16(s))
				}
				h.Write(pcmLE)
			}
			runtime.KeepAlive(pcm)
			runtime.KeepAlive(out)
			runtime.KeepAlive(packet)
			if got := fmt.Sprintf("%x", encoded.Sum(nil)); got != tc.wantPackets {
				t.Fatalf("encoder bitstream changed: got %s, want %s", got, tc.wantPackets)
			}
			digest := fmt.Sprintf("%x", h.Sum(nil))
			if digest != tc.want {
				t.Fatalf("codec output changed: got %s, want %s", digest, tc.want)
			}
		})
	}
}

// A steady 20 ms stereo frame exercises both the SILK and CELT paths.
func benchmarkPCM() []int16 {
	pcm := cBuf[int16](960 * 2)
	for i := range 960 {
		pcm[2*i] = int16(12000 * math.Sin(2*math.Pi*440*float64(i)/48000))
		pcm[2*i+1] = int16(10000 * math.Sin(2*math.Pi*660*float64(i)/48000))
	}
	return pcm
}

func BenchmarkEncode(b *testing.B) {
	tls := crtNewTLS()
	defer tls.Close()
	status := &cBuf[int32](1)[0]
	enc := opus_encoder_create(tls, 48000, 2, OPUS_APPLICATION_AUDIO, uintptr(unsafe.Pointer(status)))
	if *status != OPUS_OK || enc == 0 {
		b.Fatalf("encoder init: %d", *status)
	}
	defer opus_encoder_destroy(tls, enc)
	pcm := benchmarkPCM()
	packet := cBuf[byte](1275)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if n := opus_encode(tls, enc, uintptr(unsafe.Pointer(&pcm[0])), 960, uintptr(unsafe.Pointer(&packet[0])), int32(len(packet))); n < 0 {
			b.Fatalf("encode: %d", n)
		}
	}
	runtime.KeepAlive(pcm)
	runtime.KeepAlive(packet)
}

func BenchmarkDecode(b *testing.B) {
	tls := crtNewTLS()
	defer tls.Close()
	status := &cBuf[int32](1)[0]
	enc := opus_encoder_create(tls, 48000, 2, OPUS_APPLICATION_AUDIO, uintptr(unsafe.Pointer(status)))
	if *status != OPUS_OK || enc == 0 {
		b.Fatalf("encoder init: %d", *status)
	}
	pcm := benchmarkPCM()
	packet := cBuf[byte](1275)
	n := opus_encode(tls, enc, uintptr(unsafe.Pointer(&pcm[0])), 960, uintptr(unsafe.Pointer(&packet[0])), int32(len(packet)))
	opus_encoder_destroy(tls, enc)
	if n < 0 {
		b.Fatalf("encode fixture: %d", n)
	}
	dec := opus_decoder_create(tls, 48000, 2, uintptr(unsafe.Pointer(status)))
	if *status != OPUS_OK || dec == 0 {
		b.Fatalf("decoder init: %d", *status)
	}
	defer opus_decoder_destroy(tls, dec)
	out := cBuf[int16](5760 * 2)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if got := opus_decode(tls, dec, uintptr(unsafe.Pointer(&packet[0])), n, uintptr(unsafe.Pointer(&out[0])), 5760, 0); got != 960 {
			b.Fatalf("decode: %d", got)
		}
	}
	runtime.KeepAlive(pcm)
	runtime.KeepAlive(packet)
	runtime.KeepAlive(out)
}
