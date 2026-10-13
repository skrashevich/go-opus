// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

package opus

// stackChunk is the size of one segment of the C stack.
const crtStackChunk = 64 << 10

type crtStackSeg struct {
	base, size uintptr
}

// Mark is a position on the C stack.
type crtMark struct {
	seg int
	sp  uintptr
}

// TLS is the per-goroutine C context: it holds the stack used for C locals
// whose address is taken and for alloca. A TLS must not be used by two
// goroutines at once.
type crtTLS struct {
	segs   []crtStackSeg
	top    crtMark // segs[top.seg] holds the next free byte, top.sp
	frames []crtMark
}

func crtNewTLS() *crtTLS { return &crtTLS{top: crtMark{seg: -1}} }

// bump returns n bytes of C stack, 16-byte aligned, growing the stack by a
// new segment when the current one is full. Segments are never moved, so
// earlier allocations stay valid.
func (t *crtTLS) bump(n uintptr) uintptr {
	n = (n + 15) &^ 15
	if s := t.top.seg; s >= 0 && t.top.sp+n <= t.segs[s].base+t.segs[s].size {
		r := t.top.sp
		t.top.sp += n
		return r
	}
	s := t.top.seg + 1
	for s < len(t.segs) && t.segs[s].size < n {
		s++
	}
	if s == len(t.segs) {
		size := max(uintptr(crtStackChunk), n)
		p := crtSysAlloc(size)
		if p == 0 {
			panic("libc: out of memory for C stack")
		}
		t.segs = append(t.segs, crtStackSeg{p, size})
	}
	r := t.segs[s].base
	t.top = crtMark{s, r + n}
	return r
}

// Alloc returns n bytes of C stack. Every Alloc must be matched by Free in
// LIFO order.
func (t *crtTLS) Alloc(n int) uintptr {
	t.frames = append(t.frames, t.top)
	return t.bump(uintptr(n))
}

// Free releases the most recent Alloc.
func (t *crtTLS) Free(n int) {
	t.top = t.frames[len(t.frames)-1]
	t.frames = t.frames[:len(t.frames)-1]
}

// ArenaSave, ArenaAlloc and ArenaRestore implement alloca: memory from
// ArenaAlloc lives until ArenaRestore rewinds the stack to the saved mark.
func (t *crtTLS) ArenaSave() crtMark          { return t.top }
func (t *crtTLS) ArenaAlloc(n uint64) uintptr { return t.bump(uintptr(n)) }
func (t *crtTLS) ArenaRestore(m crtMark)      { t.top = m }

// Close releases the stack memory.
func (t *crtTLS) Close() {
	for _, s := range t.segs {
		crtSysFree(s.base)
	}
	*t = crtTLS{top: crtMark{seg: -1}}
}
