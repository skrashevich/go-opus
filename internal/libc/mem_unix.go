//go:build unix && !libc.goheap

package libc

import (
	"sync"
	"syscall"
	"unsafe"
)

// C memory lives outside the Go heap, so the uintptr arithmetic in lib.go
// passes -race/checkptr and the GC never scans or frees it.

var (
	mu     sync.Mutex
	blocks = map[uintptr][]byte{}
)

func sysAlloc(n uintptr) uintptr {
	b, err := syscall.Mmap(-1, 0, int(n), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	if err != nil {
		return 0 // malloc failure; callers report OPUS_ALLOC_FAIL
	}
	p := uintptr(unsafe.Pointer(&b[0]))
	mu.Lock()
	blocks[p] = b
	mu.Unlock()
	return p
}

func sysFree(p uintptr) {
	mu.Lock()
	b := blocks[p]
	delete(blocks, p)
	mu.Unlock()
	if b == nil {
		panic("libc: free of unknown pointer")
	}
	if err := syscall.Munmap(b); err != nil {
		panic("libc: munmap: " + err.Error())
	}
}
