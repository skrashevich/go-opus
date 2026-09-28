// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

// Package libc is the minimal C runtime that the transpiled libopus in the
// parent package needs. It replaces modernc.org/libc and implements only the
// symbols lib.go uses, with the same names and semantics.
package libc

import (
	"math"
	"math/bits"
	"unsafe"
)

// Tsize_t is the platform size_t.
type Tsize_t = uintptr

// Conversions emitted by ccgo. All are plain Go conversions.

func Bool(v bool) bool                     { return v }
func Float32FromFloat32(n float32) float32 { return n }
func Int16FromInt32(n int32) int16         { return int16(n) }
func Int16FromUint8(n uint8) int16         { return int16(n) }
func Int32FromInt64(n int64) int32         { return int32(n) }
func Int32FromUint8(n uint8) int32         { return int32(n) }
func Int32FromUint16(n uint16) int32       { return int32(n) }
func Int32FromUint32(n uint32) int32       { return int32(n) }
func Int32FromUint64(n uint64) int32       { return int32(n) }
func Uint8FromInt32(n int32) uint8         { return uint8(n) }
func Uint16FromInt16(n int16) uint16       { return uint16(n) }
func Uint16FromInt32(n int32) uint16       { return uint16(n) }
func Uint32FromInt8(n int8) uint32         { return uint32(n) }
func Uint32FromInt16(n int16) uint32       { return uint32(n) }
func Uint32FromInt32(n int32) uint32       { return uint32(n) }
func Uint32FromUint32(n uint32) uint32     { return n }
func Uint64FromInt16(n int16) uint64       { return uint64(n) }
func Uint64FromInt64(n int64) uint64       { return uint64(n) }
func Uint64FromUint64(n uint64) uint64     { return n }
func UintptrFromInt32(n int32) uintptr     { return uintptr(n) }

func BoolInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func BoolUint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

func BoolUintptr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func Xabs(t *TLS, j int32) int32 {
	if j >= 0 {
		return j
	}
	return -j
}

func Xcos(t *TLS, x float64) float64 { return math.Cos(x) }

func X__builtin_clz(t *TLS, n uint32) int32 { return int32(bits.LeadingZeros32(n)) }

func bytesAt(p uintptr, n Tsize_t) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(p)), n)
}

func Xmemcpy(t *TLS, dest, src uintptr, n Tsize_t) uintptr {
	if n != 0 {
		copy(bytesAt(dest, n), bytesAt(src, n))
	}
	return dest
}

func Xmemmove(t *TLS, dest, src uintptr, n Tsize_t) uintptr {
	return Xmemcpy(t, dest, src, n) // copy handles overlap
}

func Xmemset(t *TLS, s uintptr, c int32, n Tsize_t) uintptr {
	if n != 0 {
		b := bytesAt(s, n)
		v := byte(c)
		for i := range b {
			b[i] = v
		}
	}
	return s
}

func Xmalloc(t *TLS, n Tsize_t) uintptr {
	if n == 0 {
		n = 1
	}
	return sysAlloc(n)
}

func Xfree(t *TLS, p uintptr) {
	if p != 0 {
		sysFree(p)
	}
}
