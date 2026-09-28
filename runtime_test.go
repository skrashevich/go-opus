package opus

import (
	"testing"
	"unsafe"

	"github.com/skrashevich/go-opus/internal/libc"
)

// TestRuntimeCheckptr encodes and decodes with every buffer in C memory, so it
// runs under -race (checkptr) and exercises TLS.Alloc, malloc and VaList.
func TestRuntimeCheckptr(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	const frame, ch = 960, 2
	status := libc.Xmalloc(tls, 4)
	pcm := libc.Xmalloc(tls, frame*ch*2)
	out := libc.Xmalloc(tls, frame*ch*2)
	pkt := libc.Xmalloc(tls, 4000)
	defer func() {
		for _, p := range []uintptr{status, pcm, out, pkt} {
			libc.Xfree(tls, p)
		}
	}()
	enc := opus_encoder_create(tls, 48000, ch, OPUS_APPLICATION_AUDIO, status)
	dec := opus_decoder_create(tls, 48000, ch, status)
	if enc == 0 || dec == 0 {
		t.Fatal("create failed")
	}
	defer opus_encoder_destroy(tls, enc)
	defer opus_decoder_destroy(tls, dec)
	va := libc.Xmalloc(tls, 8)
	defer libc.Xfree(tls, va)
	if r := opus_encoder_ctl(tls, enc, OPUS_SET_BITRATE_REQUEST, libc.VaList(va, int32(64000))); r != OPUS_OK {
		t.Fatalf("ctl: %d", r)
	}
	s := unsafe.Slice((*int16)(unsafe.Pointer(pcm)), frame*ch)
	for i := 0; i < 50; i++ {
		for j := range s {
			s[j] = int16((i*977 + j*131) % 20000)
		}
		n := opus_encode(tls, enc, pcm, frame, pkt, 4000)
		if n <= 0 {
			t.Fatalf("encode: %d", n)
		}
		if d := opus_decode(tls, dec, pkt, n, out, frame, 0); d != frame {
			t.Fatalf("decode: %d", d)
		}
	}
}
