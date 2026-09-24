package opus

import (
	"crypto/sha256"
	"fmt"
	"math"
	"runtime"
	"testing"
	"unsafe"

	"modernc.org/libc"
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
	}{
		{"stereo audio", 48000, 960, 0, 2, OPUS_APPLICATION_AUDIO, 0, 10, "7100148eac64468b99fa66324ed57a454acd1d4ae257530dfc1c46e185883192"},
		{"mono voip", 48000, 960, 0, 1, OPUS_APPLICATION_VOIP, 0, 10, "ad1d8ea999f460474691cb2164cce788ef52e4433f8aca9c7fa55a28d5fb8650"},
		{"stereo noise", 48000, 960, 0, 2, OPUS_APPLICATION_AUDIO, 1, 10, "d55f39654d51200301f7e2bce6aacb712f20af0a9322e350a5f60976e38c55fd"},
		{"silk 16k mono voip", 16000, 320, 16000, 1, OPUS_APPLICATION_VOIP, 0, 10, "8f3fefa3be77c2fae500056bd583497a2ad3176664b1b099e79d51c22fb37672"},
		{"silk 8k mono noise", 8000, 160, 12000, 1, OPUS_APPLICATION_VOIP, 1, 10, "4bb1df117f6ef68400a0d759b789c4f84f1a836e9cdccb30fb261333610477e1"},
		{"hybrid stereo 24k", 48000, 960, 24000, 2, OPUS_APPLICATION_AUDIO, 0, 10, "07a5bb2972ae6f1abaa34ee1c9a2491b2cab5313695f0ab991c45addaef96b88"},
		{"celt stereo 10ms noise", 48000, 480, 96000, 2, OPUS_APPLICATION_AUDIO, 1, 10, "70f7cd94e110de21b2566ea866fd29067eff6be9143ba900a5635e14c03a676c"},
		{"mix stereo audio long", 48000, 960, 0, 2, OPUS_APPLICATION_AUDIO, 2, 200, "57d607b522c5a3e7d8debd405c97bc04963c3c97427d8c93f157ec007e8f6f23"},
		{"mix stereo 10ms 64k long", 48000, 480, 64000, 2, OPUS_APPLICATION_AUDIO, 2, 200, "bcfb5ec04c459fe04e1e1e509010be385a0302efccb0ad119f7bcc2b23fc88d8"},
		{"mix hybrid mono long", 48000, 960, 20000, 1, OPUS_APPLICATION_VOIP, 2, 200, "c0d09d637694b4dc8adf135b52705001f6d0dd77a73b0c797e6b50a9327dc2ea"},
		{"mix silk 16k long", 16000, 320, 14000, 1, OPUS_APPLICATION_VOIP, 2, 200, "4ed553920c3b613b79c7699f7627e54af1d4e577cf14ebe9faa0bc340d5b9ef4"},
		{"mix 24k stereo 40ms long", 24000, 960, 48000, 2, OPUS_APPLICATION_AUDIO, 2, 100, "322197c02cb1dd2d352d2d5a2dd1b5a80c02aaafe456e72a9169537e2372b5a0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tls := libc.NewTLS()
			defer tls.Close()
			status := &cBuf[int32](1)[0]
			enc := opus_encoder_create(tls, tc.rate, tc.channels, tc.application, uintptr(unsafe.Pointer(status)))
			if *status != OPUS_OK || enc == 0 {
				t.Fatalf("encoder init: %d", *status)
			}
			defer opus_encoder_destroy(tls, enc)
			if tc.bitrate != 0 {
				va := tls.Alloc(16)
				rc := opus_encoder_ctl(tls, enc, OPUS_SET_BITRATE_REQUEST, libc.VaList(va, tc.bitrate))
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
				got := opus_decode(tls, dec, uintptr(unsafe.Pointer(&packet[0])), n, uintptr(unsafe.Pointer(&out[0])), tc.frameSize, 0)
				if got <= 0 || got > tc.frameSize {
					t.Fatalf("decode frame %d: %d", frame, got)
				}
				h.Write(unsafe.Slice((*byte)(unsafe.Pointer(&out[0])), int(got*tc.channels*2)))
			}
			runtime.KeepAlive(pcm)
			runtime.KeepAlive(out)
			runtime.KeepAlive(packet)
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
	tls := libc.NewTLS()
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
	tls := libc.NewTLS()
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
