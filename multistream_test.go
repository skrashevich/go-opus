package opus

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"runtime"
	"slices"
	"testing"
	"unsafe"

	"github.com/skrashevich/go-opus/internal/libc"
)

// TestMultistreamOutput pins the multistream encoder and decoder output. Besides
// the packets and PCM it hashes the final range of every per-stream encoder and
// decoder, fetched through OPUS_MULTISTREAM_GET_*_STATE, which covers the
// pointer-to-pointer paths of the multistream structs.
func TestMultistreamOutput(t *testing.T) {
	for _, tc := range []struct {
		name            string
		rate, frameSize int32
		bitrate         int32
		streams         int32
		coupled         int32
		mapping         []byte
		application     int32
		frames          int
		want            string
	}{
		{"5.1 audio", 48000, 960, 256000, 4, 2, []byte{0, 4, 1, 2, 3, 5}, OPUS_APPLICATION_AUDIO, 50, "772421f453220e1460b6a6fed038856660194f8cde36a24783a7f2a72c81cc0a"},
		{"quad voip 20ms 16k", 16000, 320, 48000, 3, 1, []byte{0, 1, 2, 3}, OPUS_APPLICATION_VOIP, 50, "491b00d8ab9e3d322ce9de838bdea89d63db8da90d4f116e9654445868a39aa1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			channels := int32(len(tc.mapping))
			tls := libc.NewTLS()
			defer tls.Close()
			status := &cBuf[int32](1)[0]
			mapping := cBuf[byte](len(tc.mapping))
			copy(mapping, tc.mapping)
			enc := opus_multistream_encoder_create(tls, tc.rate, channels, tc.streams, tc.coupled, uintptr(unsafe.Pointer(&mapping[0])), tc.application, uintptr(unsafe.Pointer(status)))
			if *status != OPUS_OK || enc == 0 {
				t.Fatalf("encoder init: %d", *status)
			}
			defer opus_multistream_encoder_destroy(tls, enc)
			va := tls.Alloc(32)
			defer tls.Free(32)
			if rc := opus_multistream_encoder_ctl(tls, enc, OPUS_SET_BITRATE_REQUEST, libc.VaList(va, tc.bitrate)); rc != OPUS_OK {
				t.Fatalf("set bitrate: %d", rc)
			}
			dec := opus_multistream_decoder_create(tls, tc.rate, channels, tc.streams, tc.coupled, uintptr(unsafe.Pointer(&mapping[0])), uintptr(unsafe.Pointer(status)))
			if *status != OPUS_OK || dec == 0 {
				t.Fatalf("decoder init: %d", *status)
			}
			defer opus_multistream_decoder_destroy(tls, dec)

			pcm := cBuf[int16](int(tc.frameSize * channels))
			out := cBuf[int16](int(tc.frameSize * channels))
			packet := cBuf[byte](1275 * int(tc.streams))
			state := &cBuf[uintptr](1)[0]
			rng := &cBuf[uint32](1)[0]
			h := sha256.New()
			buf := make([]byte, 0, 2*len(out))
			// finalRanges returns the range coder state of each per-stream codec.
			finalRanges := func(st uintptr, getState int32, ctl func(*libc.TLS, uintptr, int32, uintptr) int32, streamCtl func(*libc.TLS, uintptr, int32, uintptr) int32) []uint32 {
				ranges := make([]uint32, tc.streams)
				for s := range tc.streams {
					*state = 0
					if rc := ctl(tls, st, getState, libc.VaList(va, s, uintptr(unsafe.Pointer(state)))); rc != OPUS_OK || *state == 0 {
						t.Fatalf("get state %d: %d", s, rc)
					}
					if rc := streamCtl(tls, *state, OPUS_GET_FINAL_RANGE_REQUEST, libc.VaList(va, uintptr(unsafe.Pointer(rng)))); rc != OPUS_OK {
						t.Fatalf("final range %d: %d", s, rc)
					}
					ranges[s] = *rng
					h.Write(binary.LittleEndian.AppendUint32(nil, *rng))
				}
				return ranges
			}
			seed := uint32(1)
			for frame := range tc.frames {
				for i := range tc.frameSize {
					tm := float64(int32(frame)*tc.frameSize+i) / float64(tc.rate)
					for c := range channels {
						seed = 1664525*seed + 1013904223
						v := 6000*math.Sin(2*math.Pi*float64(150+110*c)*tm)*(0.6+0.4*math.Sin(2*math.Pi*0.7*tm)) + float64(int16(seed>>16))/64
						pcm[i*channels+c] = int16(v)
					}
				}
				n := opus_multistream_encode(tls, enc, uintptr(unsafe.Pointer(&pcm[0])), tc.frameSize, uintptr(unsafe.Pointer(&packet[0])), int32(len(packet)))
				if n <= 0 {
					t.Fatalf("encode frame %d: %d", frame, n)
				}
				h.Write(packet[:n])
				encRanges := finalRanges(enc, OPUS_MULTISTREAM_GET_ENCODER_STATE_REQUEST, opus_multistream_encoder_ctl, opus_encoder_ctl)
				got := opus_multistream_decode(tls, dec, uintptr(unsafe.Pointer(&packet[0])), n, uintptr(unsafe.Pointer(&out[0])), tc.frameSize, 0)
				if got != tc.frameSize {
					t.Fatalf("decode frame %d: %d", frame, got)
				}
				// The decoder must end every stream in the encoder's range coder state.
				if decRanges := finalRanges(dec, OPUS_MULTISTREAM_GET_DECODER_STATE_REQUEST, opus_multistream_decoder_ctl, opus_decoder_ctl); !slices.Equal(encRanges, decRanges) {
					t.Fatalf("frame %d: final range mismatch: encoder %x, decoder %x", frame, encRanges, decRanges)
				}
				buf = buf[:0]
				for _, s := range out[:got*channels] {
					buf = binary.LittleEndian.AppendUint16(buf, uint16(s))
				}
				h.Write(buf)
			}
			runtime.KeepAlive(pcm)
			runtime.KeepAlive(out)
			runtime.KeepAlive(packet)
			runtime.KeepAlive(mapping)
			if digest := fmt.Sprintf("%x", h.Sum(nil)); digest != tc.want {
				t.Fatalf("multistream output changed: got %s, want %s", digest, tc.want)
			}
		})
	}
}
