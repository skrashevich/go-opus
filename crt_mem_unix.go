// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

//go:build unix && !libc.goheap

package opus

import (
	"sync"
	"syscall"
	"unsafe"
)

// C memory lives outside the Go heap, so the uintptr arithmetic in lib.go
// passes -race/checkptr and the GC never scans or frees it.

var (
	crtMu     sync.Mutex
	crtBlocks = map[uintptr][]byte{}
)

func crtSysAlloc(n uintptr) uintptr {
	b, err := syscall.Mmap(-1, 0, int(n), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	if err != nil {
		return 0 // malloc failure; callers report OPUS_ALLOC_FAIL
	}
	p := uintptr(unsafe.Pointer(&b[0]))
	crtMu.Lock()
	crtBlocks[p] = b
	crtMu.Unlock()
	return p
}

func crtSysFree(p uintptr) {
	crtMu.Lock()
	b := crtBlocks[p]
	delete(crtBlocks, p)
	crtMu.Unlock()
	if b == nil {
		panic("libc: free of unknown pointer")
	}
	if err := syscall.Munmap(b); err != nil {
		panic("libc: munmap: " + err.Error())
	}
}
