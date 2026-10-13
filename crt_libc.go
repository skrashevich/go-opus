// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

// The crt helpers provide the minimal C runtime required by the transpiled
// codec. They live in this package so Go module replacements do not depend on
// the module's canonical import path.
package opus

import (
	"math"
	"math/bits"
	"unsafe"
)

// crtTsize_t is the platform size_t.
type crtTsize_t = uintptr

// Conversions emitted by ccgo. All are plain Go conversions.

func crtBool(v bool) bool                     { return v }
func crtFloat32FromFloat32(n float32) float32 { return n }
func crtInt16FromInt32(n int32) int16         { return int16(n) }
func crtInt16FromUint8(n uint8) int16         { return int16(n) }
func crtInt32FromInt64(n int64) int32         { return int32(n) }
func crtInt32FromUint8(n uint8) int32         { return int32(n) }
func crtInt32FromUint16(n uint16) int32       { return int32(n) }
func crtInt32FromUint32(n uint32) int32       { return int32(n) }
func crtInt32FromUint64(n uint64) int32       { return int32(n) }
func crtUint8FromInt32(n int32) uint8         { return uint8(n) }
func crtUint16FromInt16(n int16) uint16       { return uint16(n) }
func crtUint16FromInt32(n int32) uint16       { return uint16(n) }
func crtUint32FromInt8(n int8) uint32         { return uint32(n) }
func crtUint32FromInt16(n int16) uint32       { return uint32(n) }
func crtUint32FromInt32(n int32) uint32       { return uint32(n) }
func crtUint32FromUint32(n uint32) uint32     { return n }
func crtUint64FromInt16(n int16) uint64       { return uint64(n) }
func crtUint64FromInt64(n int64) uint64       { return uint64(n) }
func crtUint64FromUint64(n uint64) uint64     { return n }
func crtUintptrFromInt32(n int32) uintptr     { return uintptr(n) }

func crtBoolInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func crtBoolUint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

func crtBoolUintptr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func crtXabs(t *crtTLS, j int32) int32 {
	if j >= 0 {
		return j
	}
	return -j
}

func crtXcos(t *crtTLS, x float64) float64 { return math.Cos(x) }

func crtX__builtin_clz(t *crtTLS, n uint32) int32 { return int32(bits.LeadingZeros32(n)) }

func crtBytesAt(p uintptr, n crtTsize_t) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(p)), n)
}

func crtXmemcpy(t *crtTLS, dest, src uintptr, n crtTsize_t) uintptr {
	if n != 0 {
		copy(crtBytesAt(dest, n), crtBytesAt(src, n))
	}
	return dest
}

func crtXmemmove(t *crtTLS, dest, src uintptr, n crtTsize_t) uintptr {
	return crtXmemcpy(t, dest, src, n) // copy handles overlap
}

func crtXmemset(t *crtTLS, s uintptr, c int32, n crtTsize_t) uintptr {
	if n != 0 {
		b := crtBytesAt(s, n)
		v := byte(c)
		for i := range b {
			b[i] = v
		}
	}
	return s
}

func crtXmalloc(t *crtTLS, n crtTsize_t) uintptr {
	if n == 0 {
		n = 1
	}
	return crtSysAlloc(n)
}

func crtXfree(t *crtTLS, p uintptr) {
	if p != 0 {
		crtSysFree(p)
	}
}
