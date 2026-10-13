// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

//go:build windows && !libc.goheap

package opus

import "syscall"

// C memory lives outside the Go heap; see mem_unix.go.

var (
	crtKernel32     = syscall.NewLazyDLL("kernel32.dll")
	crtVirtualAlloc = crtKernel32.NewProc("VirtualAlloc")
	crtVirtualFree  = crtKernel32.NewProc("VirtualFree")
)

const (
	crtMemCommit     = 0x1000
	crtMemReserve    = 0x2000
	crtMemRelease    = 0x8000
	crtPageReadWrite = 0x04
)

func crtSysAlloc(n uintptr) uintptr {
	p, _, _ := crtVirtualAlloc.Call(0, n, crtMemCommit|crtMemReserve, crtPageReadWrite)
	return p // 0 on failure, like malloc
}

func crtSysFree(p uintptr) {
	if r, _, err := crtVirtualFree.Call(p, 0, crtMemRelease); r == 0 {
		panic("libc: VirtualFree: " + err.Error())
	}
}
