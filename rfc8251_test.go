// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0

package opus

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
	"unsafe"
)

type rfc8251Fixture struct {
	Cases []struct {
		Name     string
		Channels int
		Packets  [][]byte
		Outputs  []struct {
			Rate, Frame    int
			FinalRanges    []uint32 `json:"final_ranges"`
			Int16SHA       string   `json:"int16_sha256"`
			FloatSHA       string   `json:"float_sha256"`
			LegacyInt16SHA string   `json:"legacy_int16_sha256"`
		}
	}
}

func loadRFC8251Fixture(t testing.TB) rfc8251Fixture {
	t.Helper()
	raw, err := os.ReadFile("testdata/rfc8251-hybrid.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture rfc8251Fixture
	if err = json.Unmarshal(raw, &fixture); err != nil {
		t.Fatal(err)
	}
	if len(fixture.Cases) != 6 {
		t.Fatalf("fixture cases: %d", len(fixture.Cases))
	}
	for _, c := range fixture.Cases {
		if c.Channels < 1 || c.Channels > 2 || len(c.Packets) == 0 {
			t.Fatalf("invalid fixture geometry: %s", c.Name)
		}
		for _, packet := range c.Packets {
			if len(packet) == 0 || len(packet) > 4000 {
				t.Fatalf("invalid fixture packet size: %s", c.Name)
			}
		}
		for _, out := range c.Outputs {
			if out.Frame < 1 || out.Frame > 5760 || len(out.FinalRanges) != len(c.Packets) {
				t.Fatalf("invalid fixture output: %s", c.Name)
			}
		}
	}
	return fixture
}

// Expectations come from the normative C source plus only RFC 8251 section 9,
// never from this Go decoder. Both APIs and entropy ranges retain their gates.
func TestRFC8251HybridReference(t *testing.T) {
	fixture := loadRFC8251Fixture(t)
	for _, c := range fixture.Cases {
		for _, want := range c.Outputs {
			for _, floating := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%d/float=%t", c.Name, want.Rate, floating), func(t *testing.T) {
					tls := NewTLS()
					defer tls.Close()
					status, pcm, packet, va := Malloc(8), Malloc(want.Frame*c.Channels*4), Malloc(4000), Malloc(16)
					defer Free(status)
					defer Free(pcm)
					defer Free(packet)
					defer Free(va)
					decoder := DecoderCreate(tls, int32(want.Rate), int32(c.Channels), status)
					if decoder == 0 || *(*int32)(unsafe.Pointer(status)) != 0 {
						t.Fatal("decoder creation")
					}
					defer DecoderDestroy(tls, decoder)
					hash := sha256.New()
					buf := make([]byte, want.Frame*c.Channels*4)
					for i, p := range c.Packets {
						copy(unsafe.Slice((*byte)(unsafe.Pointer(packet)), len(p)), p)
						var n int32
						if floating {
							n = opus_decode_float(tls, decoder, packet, int32(len(p)), pcm, int32(want.Frame), 0)
						} else {
							n = Decode(tls, decoder, packet, int32(len(p)), pcm, int32(want.Frame), 0)
						}
						if int(n) != want.Frame {
							t.Fatalf("packet %d: samples=%d", i, n)
						}
						if rc := DecoderCtl(tls, decoder, OPUS_GET_FINAL_RANGE_REQUEST, VaList(va, status)); rc != 0 {
							t.Fatal(rc)
						}
						if got := *(*uint32)(unsafe.Pointer(status)); got != want.FinalRanges[i] {
							t.Fatalf("packet %d final range: %d want %d", i, got, want.FinalRanges[i])
						}
						width := 2
						if floating {
							width = 4
							for j, v := range unsafe.Slice((*float32)(unsafe.Pointer(pcm)), want.Frame*c.Channels) {
								binary.LittleEndian.PutUint32(buf[4*j:], math.Float32bits(v))
							}
						} else {
							for j, v := range unsafe.Slice((*int16)(unsafe.Pointer(pcm)), want.Frame*c.Channels) {
								binary.LittleEndian.PutUint16(buf[2*j:], uint16(v))
							}
						}
						hash.Write(buf[:want.Frame*c.Channels*width])
					}
					expected := want.Int16SHA
					if floating {
						expected = want.FloatSHA
					}
					if got := fmt.Sprintf("%x", hash.Sum(nil)); got != expected {
						t.Fatalf("C RFC 6716 + section 9 output: got %s, want %s", got, expected)
					}
					// Section 9 preserves CELT, SILK and 16 kHz integer output for these inputs.
					if (want.Rate == 16000 || c.Name == "CELT48mono64_c10-control" || c.Name == "SILK16mono16_c10-control") && want.Int16SHA != want.LegacyInt16SHA {
						t.Fatal("fixture changes a compatibility control")
					}
				})
			}
		}
	}
}

func TestRFC8251PublicFloatReference(t *testing.T) {
	for _, c := range loadRFC8251Fixture(t).Cases {
		for _, want := range c.Outputs {
			t.Run(fmt.Sprintf("%s/%d", c.Name, want.Rate), func(t *testing.T) {
				d, err := NewDecoderWithOutput(want.Rate, c.Channels)
				if err != nil {
					t.Fatal(err)
				}
				out := make([]float32, want.Frame*c.Channels)
				buf := make([]byte, len(out)*4)
				hash := sha256.New()
				for _, p := range c.Packets {
					n, err := d.DecodeToFloat32(p, out)
					if err != nil || n != want.Frame {
						t.Fatal(n, err)
					}
					for i, v := range out {
						binary.LittleEndian.PutUint32(buf[4*i:], math.Float32bits(v))
					}
					hash.Write(buf)
				}
				if got := fmt.Sprintf("%x", hash.Sum(nil)); got != want.FloatSHA {
					t.Fatalf("public C float reference: got %s, want %s", got, want.FloatSHA)
				}
			})
		}
	}
}

func TestRFC8251HybridSteadyAllocations(t *testing.T) {
	c := loadRFC8251Fixture(t).Cases[0]
	d, err := NewDecoderWithOutput(48000, c.Channels)
	if err != nil {
		t.Fatal(err)
	}
	out := make([]float32, 960*c.Channels)
	for range 3 {
		for _, p := range c.Packets {
			if _, err := d.DecodeToFloat32(p, out); err != nil {
				t.Fatal(err)
			}
		}
	}
	i := 0
	if got := testing.AllocsPerRun(100, func() {
		n, err := d.DecodeToFloat32(c.Packets[i], out)
		if err != nil || n != 960 {
			t.Fatal(n, err)
		}
		i = (i + 1) % len(c.Packets)
	}); got != 0 {
		t.Fatalf("steady allocations: %v", got)
	}
}
