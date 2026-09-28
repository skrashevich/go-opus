// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

package opus

import (
	"runtime"
	"sync"

	"github.com/skrashevich/go-opus/internal/libc"
)

// Exported entry points of the libopus C API (opus.h) for callers outside
// this package, such as cmd/opus_demo. They follow the transpiled calling
// convention: every call takes a C context (*TLS), pointers are uintptr values
// into C memory from Malloc, and ctl requests take a VaList.
//
// The context parameter is generic so that code written against earlier
// versions, which passed a *modernc.org/libc.TLS, still compiles unchanged.
// A *TLS is used directly. For any other pointer the call borrows a TLS from
// a pool and returns it afterwards: the codec keeps no state in the TLS
// between calls, so the caller's pointer is never dereferenced or retained.
// Memory from modernc.org/libc (Xmalloc, VaList) is plain C memory and can
// be passed as before.

// TLS is the per-goroutine C context of the codec. A TLS must not be used by
// two goroutines at once.
type TLS = libc.TLS

// NewTLS returns a new C context. Call Close when done with it.
func NewTLS() *TLS { return libc.NewTLS() }

// Malloc allocates n bytes of C memory outside the Go heap, or returns 0 if
// the allocation fails. The memory must be released with Free. Each call maps
// at least one page, so allocate buffers once and reuse them.
func Malloc(n int) uintptr { return libc.Xmalloc(nil, libc.Tsize_t(n)) }

// Free releases memory returned by Malloc. Free(0) is a no-op.
func Free(p uintptr) { libc.Xfree(nil, p) }

// VaList stores ctl arguments (int32 or uintptr) at p, which must hold 8
// bytes per argument, and returns p for use as the va argument of a ctl call.
func VaList(p uintptr, args ...any) uintptr { return libc.VaList(p, args...) }

// spare holds idle contexts lent to callers that pass a foreign context.
var spare struct {
	sync.Mutex
	list []*TLS
}

// acquire returns the C context for the context argument t, and whether it
// was borrowed from spare and must be given back with release.
func acquire[T any](t *T) (*TLS, bool) {
	if own, ok := any(t).(*TLS); ok {
		return own, false
	}
	spare.Lock()
	defer spare.Unlock()
	if n := len(spare.list); n > 0 {
		own := spare.list[n-1]
		spare.list = spare.list[:n-1]
		return own, true
	}
	return libc.NewTLS(), true
}

func release(own *TLS, borrowed bool) {
	if !borrowed {
		return
	}
	spare.Lock()
	defer spare.Unlock()
	if len(spare.list) < runtime.GOMAXPROCS(0) {
		spare.list = append(spare.list, own)
		return
	}
	own.Close()
}

func EncoderCreate[T any](tls *T, fs int32, channels int32, application int32, errPtr uintptr) uintptr {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return opus_encoder_create(own, fs, channels, application, errPtr)
}

func Encode[T any](tls *T, st uintptr, pcm uintptr, frameSize int32, data uintptr, maxDataBytes int32) int32 {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return opus_encode(own, st, pcm, frameSize, data, maxDataBytes)
}

func EncoderCtl[T any](tls *T, st uintptr, request int32, va uintptr) int32 {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return opus_encoder_ctl(own, st, request, va)
}

func EncoderDestroy[T any](tls *T, st uintptr) {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	opus_encoder_destroy(own, st)
}

func DecoderCreate[T any](tls *T, fs int32, channels int32, errPtr uintptr) uintptr {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return opus_decoder_create(own, fs, channels, errPtr)
}

func Decode[T any](tls *T, st uintptr, data uintptr, length int32, pcm uintptr, frameSize int32, decodeFEC int32) int32 {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return opus_decode(own, st, data, length, pcm, frameSize, decodeFEC)
}

func DecoderCtl[T any](tls *T, st uintptr, request int32, va uintptr) int32 {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return opus_decoder_ctl(own, st, request, va)
}

func DecoderDestroy[T any](tls *T, st uintptr) {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	opus_decoder_destroy(own, st)
}

// Strerror returns a pointer to a NUL-terminated message for an error code.
func Strerror[T any](tls *T, code int32) uintptr {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return opus_strerror(own, code)
}

// GetVersionString returns a pointer to the NUL-terminated version string.
func GetVersionString[T any](tls *T) uintptr {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return opus_get_version_string(own)
}

// SIG2WORD16 and LPC_inverse_pred_gain_QA are internal helpers that earlier
// versions exported by accident; they are kept for compatibility.

func SIG2WORD16[T any](tls *T, x float32) float32 {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return sig2Word16(own, x)
}

func LPC_inverse_pred_gain_QA[T any](tls *T, A_QA uintptr, order int32) int32 {
	own, borrowed := acquire(tls)
	defer release(own, borrowed)
	return lpcInversePredGainQA(own, A_QA, order)
}
