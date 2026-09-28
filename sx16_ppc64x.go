//go:build ppc64 || ppc64le

package opus

// extBarrier is always zero. The compiler cannot prove that, so it stays in
// the expression below.
var extBarrier int32

// sx16 converts x to int16.
//
// Go's PPC64 lowering rule
//
//	(MOV(HZ|H)reg (S(R|RA)Wconst [c] x)) && x.Type.Size() <= 16 => (S(R|RA)Wconst [c] x)
//
// compares a size in bytes with 16, so it drops the 16-bit sign extension of
// any 32-bit shift by fewer than 16 bits: int32(int16(x>>6)) becomes x>>6.
// This breaks SILK, whose fixed-point macros truncate Q-format gains that way.
// XOR with a value the compiler cannot see through keeps the shift from
// feeding the extension directly. The call sites are the int16 conversions
// in lib.go where the rule fired; on other targets sx16 is int16 itself.
func sx16(x int32) int16 { return int16(x ^ extBarrier) }
