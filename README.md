# go-opus

Pure Go implementation of the [Opus audio codec](https://opus-codec.org/) (RFC 6716), transpiled from the C reference implementation using [modernc.org/ccgo v4](https://pkg.go.dev/modernc.org/ccgo/v4).

**No CGo. No system dependencies. Just Go.**

This is a complete Opus codec (SILK + CELT + Hybrid modes) — unlike other pure Go implementations that only support SILK.

## Features

- Full Opus encoder and decoder (SILK, CELT, and Hybrid modes)
- Sample rates: 8, 12, 16, 24, 48 kHz
- Bitrates: 6–510 kb/s (VBR and CBR)
- Mono and stereo
- Frame sizes: 2.5, 5, 10, 20, 40, 60 ms
- Multistream support
- Forward Error Correction (FEC)
- Zero dependencies: the module requires nothing outside the Go standard library
- Safe to use from many goroutines at once, one C context (`TLS`) per goroutine

## Performance

Apple M1 Pro, encoding and decoding a 3:41 stereo 48 kHz PCM file at 128 kb/s. Times are user CPU time, median of 5 runs, with real output files. FFmpeg runs with `-threads 1`.

| Implementation | Encode | Decode | Notes |
|---|---:|---:|---|
| C original (RFC 6716) | 1.47 s (1.0x) | 0.59 s (1.0x) | Reference implementation, `-O2` |
| **go-opus** | **1.54 s (1.05x)** | **0.75 s (1.27x)** | Pure Go, no CGo |
| FFmpeg libopus | 1.78 s (1.21x) | 0.54 s (0.92x) | SIMD-optimized |

go-opus encodes at 144x realtime and decodes at 295x realtime on this file.

The performance figures above and the following microbenchmarks were measured before checkout consolidation; they are retained measurements, not a fresh benchmark of the merged sources.

Go microbenchmarks for one 20 ms stereo 48 kHz frame (default encoder settings, codec state reused, no file I/O): encode 79 µs, decode 51 µs, with 0 allocations per frame. To reproduce:

```bash
go test -run '^$' -bench 'Benchmark(Encode|Decode)$' -benchtime=1s -count=5 -cpu=1
```

## Installation

```bash
go get github.com/skrashevich/go-opus@v0.1.0
```

### Command-line tools

```bash
# Build opus_demo (encoder/decoder)
go build -o opus_demo ./cmd/opus_demo/

# Build opus_compare (audio quality comparison)
go build -o opus_compare ./cmd/opus_compare/
```

## Usage

### Library

#### Pion-compatible Go API

The root package also exports the public `github.com/pion/opus` codec API:
`Decoder`, `NewDecoder`, `NewDecoderWithOutput`, all PCM decode methods and
`DecodePLC`; `Encoder`, `NewEncoder`, encoder options, controls and PCM encode
methods; `Bandwidth`, `Application` and their constants. `Configuration` is
retained for v0.1.0 callers. `pkg/oggreader` provides the public Ogg reader API.
Compatibility is based on Pion v0.1.0 and upstream commit
[`44637de`](https://github.com/pion/opus/tree/44637de087b3f506170ab738ed96ac1be338006b).

To switch an import directly:

```go
import opus "github.com/skrashevich/go-opus"

decoder := opus.NewDecoder() // 48 kHz mono, like pion/opus
pcm := make([]int16, 5760)   // up to 120 ms
samplesPerChannel, err := decoder.DecodeToInt16(packet, pcm)
// On success, pcm[:samplesPerChannel] contains the decoded mono samples.
```

To keep existing `github.com/pion/opus` and `github.com/pion/opus/pkg/oggreader`
imports unchanged, add one replacement to the consuming project's `go.mod`:

```go
replace github.com/pion/opus => github.com/skrashevich/go-opus v0.1.0
```

Then run `go mod tidy`. Go downloads the tagged module from GitHub through the
configured module proxy; no local checkout or adapter module is needed. The
existing `require github.com/pion/opus ...` directive stays in place. For a new
project using the Pion import path, add it with
`go mod edit -require=github.com/pion/opus@v0.1.0` before `go mod tidy`.

Alternatively, run `go get github.com/skrashevich/go-opus@v0.1.0` and change the
import to `import opus "github.com/skrashevich/go-opus"` (and, if used,
`"github.com/skrashevich/go-opus/pkg/oggreader"`). Both forms expose the same
public API. Go 1.26.1 or later is required, as declared in `go.mod`.

`NewDecoderWithOutput` supports 8, 12, 16, 24 and 48 kHz, mono or stereo.
`DecodeToInt16` and `DecodeToFloat32` return samples **per channel**; stereo PCM
is interleaved. `Decode` writes signed 16-bit little-endian bytes. Metadata
methods report the packet bandwidth and encoded stereo flag, independently of
the configured output format. Empty packets are errors; use `DecodePLC` for loss.
`Init` resets a decoder and validates settings before replacing its state.

The compatible encoder follows Pion's defaults: 48 kHz mono, 24 kbit/s,
complexity 5, CBR, constrained VBR enabled, and `ApplicationAudio`.
`Encode`/`EncodeFloat32` take one 20 ms frame and emit CELT packets;
`EncodeSILK` takes mono 20/40/60 ms PCM at 8/12/16 kHz for narrow/medium/wide
bandwidth. Like the referenced Pion encoder, the application profile is recorded
without automatically changing the other controls.

This is public API compatibility, not identical encoded bytes or bit-exact PCM
with Pion's implementation: the codec engine remains this project's RFC 6716
reference implementation. Pion's `internal` packages are outside this contract.
Go slices are copied into reusable native buffers; resources are reclaimed
automatically. Use separate encoders/decoders for independent streams, never
access an instance concurrently, and do not copy an initialized instance.
The root module still has no external dependencies.

#### Low-level API

The package exports the libopus C API (`opus.h`) in its transpiled form: pointers are `uintptr` values into C memory, and every call takes a C context.

```go
tls := opus.NewTLS() // one per goroutine
defer tls.Close()

errp := opus.Malloc(4)
pcm := opus.Malloc(960 * 2 * 2) // 20 ms of 48 kHz stereo int16
pkt := opus.Malloc(1500)
va := opus.Malloc(8) // va_list slot for ctl calls
defer func() {
	for _, p := range []uintptr{errp, pcm, pkt, va} {
		opus.Free(p)
	}
}()

enc := opus.EncoderCreate(tls, 48000, 2, opus.OPUS_APPLICATION_AUDIO, errp)
defer opus.EncoderDestroy(tls, enc)
opus.EncoderCtl(tls, enc, opus.OPUS_SET_BITRATE_REQUEST, opus.VaList(va, int32(64000)))

samples := unsafe.Slice((*int16)(unsafe.Pointer(pcm)), 960*2) // fill with audio
n := opus.Encode(tls, enc, pcm, 960, pkt, 1500)                // n bytes at pkt
```

Buffers passed to the codec must come from `opus.Malloc`, not from Go slices: C memory lives outside the Go heap, so the garbage collector never moves or frees it and `-race` (checkptr) accepts the pointer arithmetic.

**Compatibility with earlier versions.** Earlier versions required `modernc.org/libc` and took a `*libc.TLS`. The context parameter is now generic, so direct calls with a `*libc.TLS` still compile and produce the same output: for any pointer other than `*opus.TLS`, the call borrows an internal context for its duration and never touches the pointer. Memory from `libc.Xmalloc` and `libc.VaList` works as before. Two forms no longer compile: taking a function value without a type (`f := opus.Encode`; write `opus.Encode[libc.TLS]`) and passing an untyped `nil` context. New code should use `opus.NewTLS`, `opus.Malloc`, `opus.Free` and `opus.VaList`, and needs no dependency.

### Demo entry points

`opus.RunDemo(args)` retains the existing library workflow. `opus.RunReferenceDemo(args)` implements the complete reference CLI, including sweep, random frame/FEC and mode-test options; `cmd/opus_demo` uses this entry point. Both call the same codec library and return an exit status. `CPUPROFILE` is handled by the thin command wrapper.

### Command-line

```bash
# Encode raw PCM to Opus
./opus_demo -e audio 48000 2 128000 input.pcm output.opus

# Decode Opus to raw PCM
./opus_demo -d 48000 2 output.opus decoded.pcm

# Compare audio quality
./opus_compare input.pcm decoded.pcm
```

### Preparing input files

Opus demo expects raw PCM input (16-bit signed little-endian). Convert with ffmpeg:

```bash
ffmpeg -i input.wav -f s16le -acodec pcm_s16le input.pcm
```

Convert decoded PCM back to WAV:

```bash
ffmpeg -f s16le -ar 48000 -ac 2 -i decoded.pcm output.wav
```

## Project Structure

```
go-opus/
├── lib.go                  # Opus library (package opus), transpiled + hand-optimized kernels
├── export.go               # Exported API: TLS, Malloc/Free, VaList, encoder/decoder entry points
├── crt_*.go                # Private C runtime in the root package (C stack, malloc, memcpy, va_list)
├── lib_test.go             # Creation, packet round-trip tests and benchmarks
├── codec_output_test.go    # Reference codec output digests and benchmarks
├── portability_test.go     # Additional output regression coverage
├── demo.go                 # Compatible RunDemo library workflow
├── reference_demo.go       # Complete RFC 6716 demo CLI workflow
├── lib_kernels_test.go     # Optimized kernels vs. the generated reference code
├── ptr64.go                # 64-bit struct layout on 32-bit targets
├── libc_shim.go            # size_t wrappers for 32-bit targets
├── layout_test.go          # Struct layout check against testdata/layout64.txt
├── multistream_test.go     # Bit-exact multistream (5.1, quad) output test
├── export_test.go          # Exported API with *TLS and with foreign (legacy) contexts
├── runtime_test.go         # Encode/decode with C memory only; runs under -race
├── concurrency_test.go     # Parallel encoders match a serial run
├── sx16*.go                # ppc64 compiler bug workaround
├── cmd/
│   ├── opus_demo/main.go   # Thin CLI wrapper over opus.RunReferenceDemo
│   └── opus_compare/main.go # Audio comparison CLI (opus_compare.c)
└── go.mod
```

## How It Was Built

1. Source: [RFC 6716](https://www.rfc-editor.org/rfc/rfc6716) reference C implementation (130 source files, ~93K lines)
2. Each `.c` file compiled to `.o.go` using `ccgo -c` with flags: `-DUSE_ALLOCA -Drestrict= -DOPUS_BUILD`
3. Object files linked into final Go source with `ccgo`
4. The runtime from [modernc.org/libc](https://pkg.go.dev/modernc.org/libc) was replaced by the private `crt_*.go` helpers, which implement only what `lib.go` uses. `alloca` memory lives on the per-context C stack, so codec instances on different goroutines never share scratch memory
5. The CLI tools were ported to plain Go and the complete demo was moved into the library; their output is byte-identical to the transpiled versions built on darwin. The packet-loss simulation (`-loss`, `-random_fec`, `-random_framesize`) uses the `rand()` sequence of the darwin build on every platform; the transpiled Linux build used musl's `rand()` and produced different loss patterns

## Alternatives

| Library | Type | Decoder | Encoder | CGo |
|---|---|---|---|---|
| **go-opus (this)** | Pure Go | SILK, CELT, Hybrid | SILK, CELT, Hybrid | No |
| [pion/opus](https://github.com/pion/opus) | Pure Go | SILK, CELT, Hybrid | CELT (48 kHz, 20 ms), mono SILK | No |
| [hraban/opus](https://github.com/hraban/opus) | CGo wrapper | SILK, CELT, Hybrid | SILK, CELT, Hybrid | Yes |

### Comparison with pion/opus

pion/opus at commit `b10510e`, same machine and 3:41 file as above. Both libraries are called in one Go process, and the time is user CPU spent inside the encode/decode calls (median of 3 runs).

**Decoding.** The input streams were produced by the C reference encoder. Accuracy is measured against the C reference decoder output with the RFC 6716 `opus_compare` tool. A stream passes at 0% or above, and 100% means identical output.

| Stream | go-opus | pion/opus |
|---|---|---|
| CELT, 128 kb/s stereo | 0.71 s, 99.9% | 0.61 s, 99.5% |
| Hybrid, 24 kb/s stereo | 0.65 s, 99.9% | 0.47 s, fails (SNR 38.8 dB vs. reference) |
| SILK, 16 kb/s mono 16 kHz | 0.15 s, bit-exact | 0.17 s, bit-exact |

pion/opus decodes CELT and Hybrid streams faster; go-opus is closer to the reference output and is the only one of the two that passes on Hybrid.

**Encoding** 48 kHz stereo at 128 kb/s with constrained VBR and complexity 10 (pion/opus supports only 20 ms CELT frames at 48 kHz here). Both produce about 128.4 kb/s of Opus payload. They were decoded with the C reference decoder and compared with the original input:

| | go-opus | pion/opus |
|---|---|---|
| Encode time | 1.53 s | 2.87 s |
| SNR | 18.9 dB | 18.4 dB |
| `opus_compare` weighted error (lower is closer) | 0.97 | 0.45 |

go-opus encodes 1.9 times as fast. The quality measures point in different directions: go-opus has a slightly higher waveform SNR, while pion/opus keeps the per-band energies closer to the original. Neither is a full perceptual score. With the same settings, the C reference encoder (RFC 6716) scores the same as go-opus: 18.85 dB and 0.97.

## Platform Support

The code was originally transpiled on `darwin/arm64`; that identifies its provenance rather than a build restriction. It runs on 64-bit targets and on 32-bit little-endian targets. It needs only the Go standard library, so it builds for every `GOOS/GOARCH` that Go supports. In the earlier validation runs marked "bit-exact" below, `TestCodecOutput` and `TestMultistreamOutput` passed, with encoded packets and decoded PCM matching `darwin/arm64` exactly.

The consolidated sources were checked with Go 1.27.1: the library builds for all 47 targets from `go tool dist list`, while both CLI executables build with `CGO_ENABLED=0` on 42 targets. iOS and Android 386/amd64/arm executables require the toolchain's external linker. Full tests were rerun on macOS arm64/amd64 and JavaScript/WebAssembly, plus the heap runtime backend and focused `-race` checks. The table below records the earlier cross-platform runtime validation; build success alone does not establish bit-exact runtime behavior on every target.

| Platform | Build | Status | Tested with |
|---|---|---|---|
| `darwin/arm64` | ✅ | bit-exact | native |
| `darwin/amd64` | ✅ | bit-exact | Rosetta 2 |
| `linux/amd64`, `linux/arm64` | ✅ | bit-exact | Docker |
| `linux/386`, `linux/arm` (32-bit) | ✅ | bit-exact | Docker + QEMU |
| `linux/riscv64` | ✅ | bit-exact | Docker + QEMU |
| `linux/s390x` (big-endian) | ✅ | bit-exact | Docker + QEMU |
| `linux/ppc64le` | ✅ | bit-exact (POWER8, POWER9; POWER10 not tested) | Docker + QEMU |
| `windows/amd64`, `windows/386` | ✅ | bit-exact | Wine |
| `windows/arm64`, `freebsd/*`, `netbsd/*`, `openbsd/*`, `illumos/amd64`, `linux/loong64` | ✅ | builds, not run | — |
| `js/wasm` | ✅ | bit-exact | Node.js (`go_js_wasm_exec`) |
| `linux/mips64` (big-endian) | ✅ | bit-exact | Docker + QEMU |
| `wasip1/wasm`, `linux/mips64le`, `linux/ppc64`, `linux/mips`, `plan9/*`, `aix/ppc64`, `solaris/amd64`, `freebsd/riscv64` | ✅ | builds, not run | — |

On targets without `mmap`/`VirtualAlloc` (wasm, plan9), C memory is Go heap memory kept alive by the runtime (`crt_mem_goheap.go`); `-tags libc.goheap` selects this backend anywhere for testing.

### Portability notes

- **32-bit targets** keep the 64-bit memory layout that the transpiled code hardcodes. Each C pointer field is padded to 8 bytes and aligned to 8 bytes (`ptr64.go`). `size_t` arguments go through small wrappers (`libc_shim.go`). `TestStructLayout` checks every codec struct against `testdata/layout64.txt`. `TestStructLayoutCoverage` fails when code starts using a struct that is not checked, or views C memory as a Go `[N]uintptr` array.
- **ppc64/ppc64le** needs a workaround for a Go compiler bug (Go 1.27.1). A lowering rule in `PPC64.rules` compares a size in bytes with 16. As a result it drops the 16-bit extension after a 32-bit shift, so `int32(int16(x >> 6))` compiles to `x >> 6`. The affected conversions in `lib.go` use `sx16` (`sx16_ppc64x.go`). On other targets `sx16` is plain `int16`.
- **Big-endian:** test digests hash PCM as little-endian, so the same digests hold on every target.

## License

Original additions and modifications to go-opus are Copyright (c) 2026
Sergey Krashevich and licensed under the [Apache License 2.0](LICENSE).
See [NOTICE](NOTICE) for project attribution.

The underlying Opus reference code and derived portions retain their original
BSD terms and copyright notices. Those notices are preserved in the source
files and [COPYING](COPYING); Apache-2.0 does not replace their requirements.
This is not a choice between Apache and BSD for the entire project: comply
with the terms applicable to the portions you redistribute.

When redistributing this version or derivatives containing its Apache-licensed
contributions, include LICENSE, preserve the applicable attribution from NOTICE
as required by Apache-2.0 section 4, and retain the applicable upstream notices
from COPYING and the source files. Mark files you modify as required by the
license. This licensing change does not revoke permissions granted for earlier
versions.

## Credits

- [Opus Codec](https://opus-codec.org/) — IETF RFC 6716
- [modernc.org/ccgo](https://pkg.go.dev/modernc.org/ccgo/v4) — C to Go transpiler by Jan Mercl
- [modernc.org/libc](https://pkg.go.dev/modernc.org/libc) — Go libc runtime used by earlier versions
