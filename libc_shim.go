package opus

import "modernc.org/libc"

// lib.go was generated for a 64-bit target, where size_t is uint64. These
// wrappers keep that signature and convert to the platform's libc.Tsize_t, so
// the same code builds on 32-bit targets. On 64-bit they inline to the direct
// libc call.

func xmemcpy(tls *libc.TLS, dest, src uintptr, n uint64) uintptr {
	return libc.Xmemcpy(tls, dest, src, libc.Tsize_t(n))
}

func xmemmove(tls *libc.TLS, dest, src uintptr, n uint64) uintptr {
	return libc.Xmemmove(tls, dest, src, libc.Tsize_t(n))
}

func xmemset(tls *libc.TLS, s uintptr, c int32, n uint64) uintptr {
	return libc.Xmemset(tls, s, c, libc.Tsize_t(n))
}

func xmalloc(tls *libc.TLS, n uint64) uintptr {
	return libc.Xmalloc(tls, libc.Tsize_t(n))
}
