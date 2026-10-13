// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

package opus

// lib.go was generated for a 64-bit target, where size_t is uint64. These
// wrappers keep that signature and convert to the platform's crtTsize_t, so
// the same code builds on 32-bit targets. On 64-bit they inline to the direct
// libc call.

func xmemcpy(tls *crtTLS, dest, src uintptr, n uint64) uintptr {
	return crtXmemcpy(tls, dest, src, crtTsize_t(n))
}

func xmemmove(tls *crtTLS, dest, src uintptr, n uint64) uintptr {
	return crtXmemmove(tls, dest, src, crtTsize_t(n))
}

func xmemset(tls *crtTLS, s uintptr, c int32, n uint64) uintptr {
	return crtXmemset(tls, s, c, crtTsize_t(n))
}

func xmalloc(tls *crtTLS, n uint64) uintptr {
	return crtXmalloc(tls, crtTsize_t(n))
}
