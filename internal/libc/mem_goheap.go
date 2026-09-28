// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

//go:build (!unix && !windows) || libc.goheap

package libc

import (
	"sync"
	"unsafe"
)

// Fallback for targets without mmap (wasm, plan9): C memory is Go heap memory
// kept alive by the blocks map. The Go GC does not move heap objects, so the
// addresses stay valid. This fallback is incompatible with checkptr (-race),
// which none of these targets support.

var (
	mu     sync.Mutex
	blocks = map[uintptr][]uint64{}
)

func sysAlloc(n uintptr) uintptr {
	b := make([]uint64, (n+7)/8) // 8-byte aligned, zeroed
	p := uintptr(unsafe.Pointer(&b[0]))
	mu.Lock()
	blocks[p] = b
	mu.Unlock()
	return p
}

func sysFree(p uintptr) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := blocks[p]; !ok {
		panic("libc: free of unknown pointer")
	}
	delete(blocks, p)
}
