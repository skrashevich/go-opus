package opus

import (
	"runtime"
	"sync"
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

	// Any other pointer type is accepted as a context; the results must be
	// the same.
	if got := roundTrip(t, new(foreignCtx)); string(got) != string(want) {
		t.Fatal("output with a foreign context differs")
	}
	if got := roundTrip(t, (*foreignCtx)(nil)); string(got) != string(want) {
		t.Fatal("output with a nil foreign context differs")
	}
}

// TestForeignContextChurn uses a fresh foreign context for every codec
// instance on several goroutines while the GC runs, so freed context
// addresses are reused. Borrowed contexts must never be shared or released
// while a call is using them.
func TestForeignContextChurn(t *testing.T) {
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				runtime.GC()
				time.Sleep(time.Millisecond)
			}
		}
	}()
	defer close(stop)

	const frame, ch = 960, 2
	var wg sync.WaitGroup
	deadline := time.Now().Add(2 * time.Second)
	for g := 0; g < max(4, runtime.GOMAXPROCS(0)); g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errp := Malloc(4)
			pcm := Malloc(frame * ch * 2)
			pkt := Malloc(1500)
			defer Free(errp)
			defer Free(pcm)
			defer Free(pkt)
			s := unsafe.Slice((*int16)(unsafe.Pointer(pcm)), frame*ch)
			for i := 0; time.Now().Before(deadline); i++ {
				enc := EncoderCreate(new(foreignCtx), 48000, ch, OPUS_APPLICATION_AUDIO, errp)
				if enc == 0 {
					t.Error("create failed")
					return
				}
				for f := 0; f < 5; f++ {
					for j := range s {
						s[j] = int16((i*977 + j*131 + f) % 20000)
					}
					if n := Encode(new(foreignCtx), enc, pcm, frame, pkt, 1500); n <= 0 {
						t.Errorf("encode: %d", n)
						return
					}
				}
				EncoderDestroy(new(foreignCtx), enc)
			}
		}()
	}
	wg.Wait()
}
