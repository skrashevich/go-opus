// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

//go:build windows && !libc.goheap

package libc

import "syscall"

// C memory lives outside the Go heap; see mem_unix.go.

var (
	kernel32     = syscall.NewLazyDLL("kernel32.dll")
	virtualAlloc = kernel32.NewProc("VirtualAlloc")
	virtualFree  = kernel32.NewProc("VirtualFree")
)

const (
	memCommit     = 0x1000
	memReserve    = 0x2000
	memRelease    = 0x8000
	pageReadWrite = 0x04
)

func sysAlloc(n uintptr) uintptr {
	p, _, _ := virtualAlloc.Call(0, n, memCommit|memReserve, pageReadWrite)
	return p // 0 on failure, like malloc
}

func sysFree(p uintptr) {
	if r, _, err := virtualFree.Call(p, 0, memRelease); r == 0 {
		panic("libc: VirtualFree: " + err.Error())
	}
}
