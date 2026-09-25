package opus

import (
	"sync/atomic"
	"unsafe"
)

// lib.go was generated for a 64-bit target: every C pointer occupies 8 bytes,
// and struct offsets, struct sizes and pointer-array strides are hardcoded
// with that assumption. On 32-bit targets the code keeps the same memory
// layout by storing each pointer in the low (little-endian) half of an 8-byte
// slot. These types are zero-sized or single-word on 64-bit targets.

// align8 precedes every pointer field in a codec struct. It takes no space but
// aligns the field, and with it the whole struct, to 8 bytes as on 64-bit
// targets: the compiler gives atomic.Int64 8-byte alignment everywhere.
type align8 [0]atomic.Int64

// ptrpad follows every uintptr field in a codec struct, widening it to 8 bytes.
type ptrpad [8/unsafe.Sizeof(uintptr(0)) - 1]uintptr

// ptrslot is an array element holding one pointer in 8 bytes; element [0] is
// the pointer.
type ptrslot [8 / unsafe.Sizeof(uintptr(0))]uintptr

// pslot returns the i-th pointer slot of a C array of pointers at p.
func pslot(p uintptr, i int32) *uintptr {
	return (*uintptr)(unsafe.Pointer(p + uintptr(i)*8))
}
