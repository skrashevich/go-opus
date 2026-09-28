package opus

import (
	"runtime"
	"testing"
	"time"
	"unsafe"
)

func cString(p uintptr) string {
	n := 0
	for *(*byte)(unsafe.Pointer(p + uintptr(n))) != 0 {
		n++
	}
	return string(unsafe.Slice((*byte)(unsafe.Pointer(p)), n))
}

// foreignCtx stands in for a caller's own context type, such as
// *modernc.org/libc.TLS in code written against earlier versions.
type foreignCtx struct{ _ [64]byte }

func foreignCount() int {
	n := 0
	foreignTLS.Range(func(_, _ any) bool { n++; return true })
	return n
}

// roundTrip encodes and decodes a few frames through the exported API.
func roundTrip[T any](t *testing.T, tls *T) []byte {
	const frame, ch = 960, 2
	errp := Malloc(4)
	pcm := Malloc(frame * ch * 2)
	out := Malloc(frame * ch * 2)
	pkt := Malloc(1500)
	va := Malloc(8)
	defer func() {
		for _, p := range []uintptr{errp, pcm, out, pkt, va} {
			Free(p)
		}
	}()
	enc := EncoderCreate(tls, 48000, ch, OPUS_APPLICATION_AUDIO, errp)
	dec := DecoderCreate(tls, 48000, ch, errp)
	if enc == 0 || dec == 0 {
		t.Fatal("create failed")
	}
	defer EncoderDestroy(tls, enc)
	defer DecoderDestroy(tls, dec)
	if r := EncoderCtl(tls, enc, OPUS_SET_BITRATE_REQUEST, VaList(va, int32(48000))); r != OPUS_OK {
		t.Fatalf("encoder ctl: %s", cString(Strerror(tls, r)))
	}
	s := unsafe.Slice((*int16)(unsafe.Pointer(pcm)), frame*ch)
	var res []byte
	for i := 0; i < 20; i++ {
		for j := range s {
			s[j] = int16((i*977 + j*131) % 20000)
		}
		n := Encode(tls, enc, pcm, frame, pkt, 1500)
		if n <= 0 {
			t.Fatalf("encode: %d", n)
		}
		if d := Decode(tls, dec, pkt, n, out, frame, 0); d != frame {
			t.Fatalf("decode: %d", d)
		}
		if r := DecoderCtl(tls, dec, OPUS_GET_BANDWIDTH_REQUEST, VaList(va, errp)); r != OPUS_OK {
			t.Fatalf("decoder ctl: %s", cString(Strerror(tls, r)))
		}
		res = append(res, unsafe.Slice((*byte)(unsafe.Pointer(errp)), 4)...)
		res = append(res, unsafe.Slice((*byte)(unsafe.Pointer(pkt)), n)...)
		res = append(res, unsafe.Slice((*byte)(unsafe.Pointer(out)), frame*ch*2)...)
	}
	return res
}

func TestExportedAPIContexts(t *testing.T) {
	tls := NewTLS()
	defer tls.Close()
	if v := cString(GetVersionString(tls)); v == "" {
		t.Fatal("empty version string")
	}
	want := roundTrip(t, tls)
	if foreignCount() != 0 {
		t.Fatal("a *TLS must be used directly, not mapped")
	}

	// Any other pointer type is accepted as a context key and mapped to an
	// internal TLS; the results must be the same.
	func() {
		key := new(foreignCtx)
		if got := roundTrip(t, key); string(got) != string(want) {
			t.Fatal("output with a foreign context differs")
		}
		if foreignCount() != 1 {
			t.Fatalf("foreign contexts mapped: %d, want 1", foreignCount())
		}
		runtime.KeepAlive(key)
	}()

	// The mapped TLS is released once the key is garbage collected.
	for deadline := time.Now().Add(5 * time.Second); foreignCount() != 0; {
		if time.Now().After(deadline) {
			t.Fatal("foreign context was not released after GC")
		}
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
}
