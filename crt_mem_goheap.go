// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

//go:build (!unix && !windows) || libc.goheap

package opus

import (
	"sync"
	"unsafe"
)

// Fallback for targets without mmap (wasm, plan9): C memory is Go heap memory
// kept alive by the blocks map. The Go GC does not move heap objects, so the
// addresses stay valid. This fallback is incompatible with checkptr (-race),
// which none of these targets support.

var (
	crtMu     sync.Mutex
	crtBlocks = map[uintptr][]uint64{}
)

func crtSysAlloc(n uintptr) uintptr {
	b := make([]uint64, (n+7)/8) // 8-byte aligned, zeroed
	p := uintptr(unsafe.Pointer(&b[0]))
	crtMu.Lock()
	crtBlocks[p] = b
	crtMu.Unlock()
	return p
}

func crtSysFree(p uintptr) {
	crtMu.Lock()
	defer crtMu.Unlock()
	if _, ok := crtBlocks[p]; !ok {
		panic("libc: free of unknown pointer")
	}
	delete(crtBlocks, p)
}
