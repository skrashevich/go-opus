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

func TestCodecOutput(t *testing.T) {
	for _, tc := range []struct {
		name        string
		channels    int32
		application int32
		noise       bool
		want        string
	}{
			{"stereo audio", 2, OPUS_APPLICATION_AUDIO, false, "dc70c435efb4ec551db3d252587c15a8648b80b6dd8ba4f3062ca767ce812f5d"},
			{"mono voip", 1, OPUS_APPLICATION_VOIP, false, "c9506875459f376010e239e6d6c16ae83169512fedd17855efa7ec5be92042c3"},
			{"stereo noise", 2, OPUS_APPLICATION_AUDIO, true, "722ac286171a4fbe4a8b45ce83557df6017a6437d18e8df5bad0b03b4eadf082"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tls := libc.NewTLS()
			defer tls.Close()
			var status int32
			enc := opus_encoder_create(tls, 48000, tc.channels, tc.application, uintptr(unsafe.Pointer(&status)))
			if status != OPUS_OK || enc == 0 {
				t.Fatalf("encoder init: %d", status)
			}
			defer opus_encoder_destroy(tls, enc)
			dec := opus_decoder_create(tls, 48000, tc.channels, uintptr(unsafe.Pointer(&status)))
			if status != OPUS_OK || dec == 0 {
				t.Fatalf("decoder init: %d", status)
			}
			defer opus_decoder_destroy(tls, dec)
			pcm := make([]int16, 960*tc.channels)
			out := make([]int16, 960*tc.channels)
			packet := make([]byte, 1275)
			h := sha256.New()
			seed := uint32(1)
			for frame := range 10 {
				for i := range pcm {
					if tc.noise {
						seed = 1664525*seed + 1013904223
						pcm[i] = int16(seed>>16) / 2
					} else {
						pcm[i] = int16(10000 * math.Sin(2*math.Pi*440*float64(frame*len(pcm)+i)/48000))
					}
				}
				n := opus_encode(tls, enc, uintptr(unsafe.Pointer(&pcm[0])), 960, uintptr(unsafe.Pointer(&packet[0])), int32(len(packet)))
				if n <= 0 {
					t.Fatalf("encode frame %d: %d", frame, n)
				}
				h.Write(packet[:n])
				got := opus_decode(tls, dec, uintptr(unsafe.Pointer(&packet[0])), n, uintptr(unsafe.Pointer(&out[0])), 960, 0)
				if got <= 0 || got > 960 {
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
	pcm := make([]int16, 960*2)
	for i := range 960 {
		pcm[2*i] = int16(12000 * math.Sin(2*math.Pi*440*float64(i)/48000))
		pcm[2*i+1] = int16(10000 * math.Sin(2*math.Pi*660*float64(i)/48000))
	}
	return pcm
}

func BenchmarkEncode(b *testing.B) {
	tls := libc.NewTLS()
	defer tls.Close()
	var status int32
	enc := opus_encoder_create(tls, 48000, 2, OPUS_APPLICATION_AUDIO, uintptr(unsafe.Pointer(&status)))
	if status != OPUS_OK || enc == 0 {
		b.Fatalf("encoder init: %d", status)
	}
	defer opus_encoder_destroy(tls, enc)
	pcm := benchmarkPCM()
	packet := make([]byte, 1275)
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
	var status int32
	enc := opus_encoder_create(tls, 48000, 2, OPUS_APPLICATION_AUDIO, uintptr(unsafe.Pointer(&status)))
	if status != OPUS_OK || enc == 0 {
		b.Fatalf("encoder init: %d", status)
	}
	pcm := benchmarkPCM()
	packet := make([]byte, 1275)
	n := opus_encode(tls, enc, uintptr(unsafe.Pointer(&pcm[0])), 960, uintptr(unsafe.Pointer(&packet[0])), int32(len(packet)))
	opus_encoder_destroy(tls, enc)
	if n < 0 {
		b.Fatalf("encode fixture: %d", n)
	}
	dec := opus_decoder_create(tls, 48000, 2, uintptr(unsafe.Pointer(&status)))
	if status != OPUS_OK || dec == 0 {
		b.Fatalf("decoder init: %d", status)
	}
	defer opus_decoder_destroy(tls, dec)
	out := make([]int16, 5760*2)
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
