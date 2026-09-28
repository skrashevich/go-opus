package opus

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/skrashevich/go-opus/internal/libc"
)

// Exported entry points of the libopus C API (opus.h) for callers outside
// this package, such as cmd/opus_demo. They follow the transpiled calling
// convention: every call takes a C context (*TLS), pointers are uintptr values
// into C memory from Malloc, and ctl requests take a VaList.
//
// The context parameter is generic so that code written against earlier
// versions, which passed a *modernc.org/libc.TLS, still compiles unchanged.
// A *TLS is used directly; any other non-nil pointer is treated as an opaque
// key and mapped to a *TLS owned by this package, which is released when the
// key is garbage collected. Memory from modernc.org/libc (Xmalloc, VaList) is
// plain C memory and can be passed as before.

// TLS is the per-goroutine C context of the codec. A TLS must not be used by
// two goroutines at once.
type TLS = libc.TLS

// NewTLS returns a new C context. Call Close when done with it.
func NewTLS() *TLS { return libc.NewTLS() }

// Malloc allocates n bytes of C memory outside the Go heap, or returns 0 if
// the allocation fails. The memory must be released with Free.
func Malloc(n int) uintptr { return libc.Xmalloc(nil, libc.Tsize_t(n)) }

// Free releases memory returned by Malloc. Free(0) is a no-op.
func Free(p uintptr) { libc.Xfree(nil, p) }

// VaList stores ctl arguments (int32 or uintptr) at p, which must hold 8
// bytes per argument, and returns p for use as the va argument of a ctl call.
func VaList(p uintptr, args ...any) uintptr { return libc.VaList(p, args...) }

var foreignTLS sync.Map // uintptr(key) -> *TLS

// tlsOf returns the C context for the context argument t.
func tlsOf[T any](t *T) *TLS {
	if own, ok := any(t).(*TLS); ok {
		return own
	}
	if t == nil {
		panic("opus: nil TLS")
	}
	k := uintptr(unsafe.Pointer(t))
	if v, ok := foreignTLS.Load(k); ok {
		return v.(*TLS)
	}
	own := libc.NewTLS()
	if v, loaded := foreignTLS.LoadOrStore(k, own); loaded {
		return v.(*TLS)
	}
	runtime.AddCleanup(t, func(k uintptr) {
		if v, ok := foreignTLS.LoadAndDelete(k); ok {
			v.(*TLS).Close()
		}
	}, k)
	return own
}

func EncoderCreate[T any](tls *T, fs int32, channels int32, application int32, errPtr uintptr) uintptr {
	return opus_encoder_create(tlsOf(tls), fs, channels, application, errPtr)
}

func Encode[T any](tls *T, st uintptr, pcm uintptr, frameSize int32, data uintptr, maxDataBytes int32) int32 {
	return opus_encode(tlsOf(tls), st, pcm, frameSize, data, maxDataBytes)
}

func EncoderCtl[T any](tls *T, st uintptr, request int32, va uintptr) int32 {
	return opus_encoder_ctl(tlsOf(tls), st, request, va)
}

func EncoderDestroy[T any](tls *T, st uintptr) {
	opus_encoder_destroy(tlsOf(tls), st)
}

func DecoderCreate[T any](tls *T, fs int32, channels int32, errPtr uintptr) uintptr {
	return opus_decoder_create(tlsOf(tls), fs, channels, errPtr)
}

func Decode[T any](tls *T, st uintptr, data uintptr, length int32, pcm uintptr, frameSize int32, decodeFEC int32) int32 {
	return opus_decode(tlsOf(tls), st, data, length, pcm, frameSize, decodeFEC)
}

func DecoderCtl[T any](tls *T, st uintptr, request int32, va uintptr) int32 {
	return opus_decoder_ctl(tlsOf(tls), st, request, va)
}

func DecoderDestroy[T any](tls *T, st uintptr) {
	opus_decoder_destroy(tlsOf(tls), st)
}

// Strerror returns a pointer to a NUL-terminated message for an error code.
func Strerror[T any](tls *T, code int32) uintptr {
	return opus_strerror(tlsOf(tls), code)
}

// GetVersionString returns a pointer to the NUL-terminated version string.
func GetVersionString[T any](tls *T) uintptr {
	return opus_get_version_string(tlsOf(tls))
}
