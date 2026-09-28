package libc

// stackChunk is the size of one segment of the C stack.
const stackChunk = 64 << 10

type stackSeg struct {
	base, size uintptr
}

// Mark is a position on the C stack.
type Mark struct {
	seg int
	sp  uintptr
}

// TLS is the per-goroutine C context: it holds the stack used for C locals
// whose address is taken and for alloca. A TLS must not be used by two
// goroutines at once.
type TLS struct {
	segs   []stackSeg
	top    Mark // segs[top.seg] holds the next free byte, top.sp
	frames []Mark
}

func NewTLS() *TLS { return &TLS{top: Mark{seg: -1}} }

// bump returns n bytes of C stack, 16-byte aligned, growing the stack by a
// new segment when the current one is full. Segments are never moved, so
// earlier allocations stay valid.
func (t *TLS) bump(n uintptr) uintptr {
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
		size := max(uintptr(stackChunk), n)
		p := sysAlloc(size)
		if p == 0 {
			panic("libc: out of memory for C stack")
		}
		t.segs = append(t.segs, stackSeg{p, size})
	}
	r := t.segs[s].base
	t.top = Mark{s, r + n}
	return r
}

// Alloc returns n bytes of C stack. Every Alloc must be matched by Free in
// LIFO order.
func (t *TLS) Alloc(n int) uintptr {
	t.frames = append(t.frames, t.top)
	return t.bump(uintptr(n))
}

// Free releases the most recent Alloc.
func (t *TLS) Free(n int) {
	t.top = t.frames[len(t.frames)-1]
	t.frames = t.frames[:len(t.frames)-1]
}

// ArenaSave, ArenaAlloc and ArenaRestore implement alloca: memory from
// ArenaAlloc lives until ArenaRestore rewinds the stack to the saved mark.
func (t *TLS) ArenaSave() Mark             { return t.top }
func (t *TLS) ArenaAlloc(n uint64) uintptr { return t.bump(uintptr(n)) }
func (t *TLS) ArenaRestore(m Mark)         { t.top = m }

// Close releases the stack memory.
func (t *TLS) Close() {
	for _, s := range t.segs {
		sysFree(s.base)
	}
	*t = TLS{top: Mark{seg: -1}}
}
