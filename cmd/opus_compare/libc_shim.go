package main

import (
	"unsafe"

	"modernc.org/libc"
)

// ptrSize is the stride of C pointer arrays such as argv, which libc builds
// with native pointers.
const ptrSize = unsafe.Sizeof(uintptr(0))

// This file was generated for a 64-bit target, where size_t is uint64. These
// wrappers keep that signature and convert to the platform's libc.Tsize_t, so
// the same code builds on 32-bit targets.

func xmalloc(tls *libc.TLS, n uint64) uintptr {
	return libc.Xmalloc(tls, libc.Tsize_t(n))
}

func xrealloc(tls *libc.TLS, p uintptr, n uint64) uintptr {
	return libc.Xrealloc(tls, p, libc.Tsize_t(n))
}

func xfread(tls *libc.TLS, ptr uintptr, size, n uint64, f uintptr) uint64 {
	return uint64(libc.Xfread(tls, ptr, libc.Tsize_t(size), libc.Tsize_t(n), f))
}
