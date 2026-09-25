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
- Zero external dependencies (pure Go via modernc.org/libc)

## Performance

Apple M1 Pro, encoding and decoding a 3:41 stereo 48 kHz PCM file at 128 kb/s. Times are user CPU time, median of 5 runs, with real output files. FFmpeg runs with `-threads 1`.

| Implementation | Encode | Decode | Notes |
|---|---:|---:|---|
| C original (RFC 6716) | 1.50 s (1.0x) | 0.61 s (1.0x) | Reference implementation, `-O2` |
| **go-opus** | **1.58 s (1.05x)** | **0.77 s (1.26x)** | Pure Go, no CGo |
| FFmpeg libopus | 1.83 s (1.22x) | 0.57 s (0.93x) | SIMD-optimized |

go-opus encodes at 140x realtime and decodes at 288x realtime on this file.

Go microbenchmarks for one 20 ms stereo 48 kHz frame (default encoder settings, codec state reused, no file I/O): encode 79 µs, decode 48 µs, with 0 allocations per frame. To reproduce:

```bash
go test -run '^$' -bench 'Benchmark(Encode|Decode)$' -benchtime=1s -count=5 -cpu=1
```

## Installation

```bash
go get github.com/skrashevich/go-opus
```

### Command-line tools

```bash
# Build opus_demo (encoder/decoder)
go build -o opus_demo ./cmd/opus_demo/

# Build opus_compare (audio quality comparison)
go build -o opus_compare ./cmd/opus_compare/
```

## Usage

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
├── export.go               # Exported encoder/decoder entry points used by the CLI
├── lib_test.go             # Bit-exact codec output tests and benchmarks
├── lib_kernels_test.go     # Optimized kernels vs. the generated reference code
├── ptr64.go                # 64-bit struct layout on 32-bit targets
├── libc_shim.go            # size_t wrappers for 32-bit targets
├── layout_test.go          # Struct layout check against testdata/layout64.txt
├── cmd/
│   ├── opus_demo/main.go   # Encoder/decoder CLI (opus_demo.c) on top of package opus
│   └── opus_compare/main.go # Audio comparison CLI
├── go.mod
└── go.sum
```

## How It Was Built

1. Source: [RFC 6716](https://www.rfc-editor.org/rfc/rfc6716) reference C implementation (130 source files, ~93K lines)
2. Each `.c` file compiled to `.o.go` using `ccgo -c` with flags: `-DUSE_ALLOCA -Drestrict= -DOPUS_BUILD`
3. Object files linked into final Go source with `ccgo`
4. Runtime provided by [modernc.org/libc](https://pkg.go.dev/modernc.org/libc)

## Alternatives

| Library | Type | Decoder | Encoder | CGo |
|---|---|---|---|---|
| **go-opus (this)** | Pure Go | SILK, CELT, Hybrid | SILK, CELT, Hybrid | No |
| [pion/opus](https://github.com/pion/opus) | Pure Go | SILK, CELT, Hybrid | CELT (48 kHz, 20 ms), mono SILK | No |
| [hraban/opus](https://github.com/hraban/opus) | CGo wrapper | SILK, CELT, Hybrid | SILK, CELT, Hybrid | Yes |

### Comparison with pion/opus

pion/opus at commit `86ced73`, same machine and 3:41 file as above. Both libraries are called in one Go process, and the time is user CPU spent inside the encode/decode calls (median of 3 runs).

**Decoding.** The input streams were produced by the C reference encoder. Accuracy is measured against the C reference decoder output with the RFC 6716 `opus_compare` tool. A stream passes at 0% or above, and 100% means identical output.

| Stream | go-opus | pion/opus |
|---|---|---|
| CELT, 128 kb/s stereo | 0.72 s, 99.9% | 0.78 s, 99.5% |
| Hybrid, 24 kb/s stereo | 0.78 s, 99.9% | 0.50 s, fails (SNR 38.8 dB vs. reference) |
| SILK, 16 kb/s mono 16 kHz | 0.16 s, bit-exact | 0.18 s, bit-exact |

**Encoding** 48 kHz stereo at 128 kb/s with constrained VBR and complexity 10 (pion/opus supports only 20 ms CELT frames at 48 kHz here). Both produce about 128.4 kb/s of Opus payload. They were decoded with the C reference decoder and compared with the original input:

| | go-opus | pion/opus |
|---|---|---|
| Encode time | 1.55 s | 3.17 s |
| SNR | 18.9 dB | 18.4 dB |
| `opus_compare` weighted error (lower is closer) | 0.97 | 0.45 |

go-opus encodes about twice as fast. The quality measures point in different directions: go-opus has a slightly higher waveform SNR, while pion/opus keeps the per-band energies closer to the original. Neither is a full perceptual score. With the same settings, the C reference encoder (RFC 6716) scores the same as go-opus: 18.85 dB and 0.97.

## Platform Support

The code was transpiled on `darwin/arm64`, but it runs on other 64-bit targets and on 32-bit little-endian targets. On every platform marked "bit-exact" below, `TestCodecOutput` passes, so encoded packets and decoded PCM match `darwin/arm64` exactly.

| Platform | Build | Status | Tested with |
|---|---|---|---|
| `darwin/arm64` | ✅ | bit-exact | native |
| `darwin/amd64` | ✅ | bit-exact | Rosetta 2 |
| `linux/amd64`, `linux/arm64` | ✅ | bit-exact | Docker |
| `linux/riscv64` | ✅ | bit-exact | Docker + QEMU |
| `linux/s390x` (big-endian) | ✅ | bit-exact | Docker + QEMU |
| `windows/amd64` | ✅ | bit-exact | Wine |
| `linux/ppc64le` | ✅ | **broken**: SILK/hybrid output is wrong in optimized builds (correct with `-gcflags='-N -l'`); CELT is fine | Docker + QEMU |
| `windows/arm64`, `freebsd/amd64`, `linux/loong64` | ✅ | builds, not run | — |
| `linux/386`, `linux/arm` (32-bit) | ✅ | bit-exact | Docker + QEMU |
| `windows/386` | ✅ | bit-exact | Wine |
| `wasip1/wasm`, `js/wasm` | ❌ | not supported by `modernc.org/libc` | — |

On Windows, `cmd/opus_demo` does not build because `modernc.org/libc` has no `feof` there. The library itself works.

32-bit targets keep the 64-bit memory layout: each C pointer field is padded to 8 bytes and 8-byte aligned (`ptr64.go`), and `size_t` arguments go through small wrappers (`libc_shim.go`). `TestStructLayout` checks the layout of every codec struct against `testdata/layout64.txt`. This only works on little-endian targets, so 32-bit big-endian targets (MIPS, PowerPC) are not supported.

## License

BSD 3-Clause — same as the original Opus reference implementation (IETF Trust, Skype Limited, Xiph.Org Foundation).

## Credits

- [Opus Codec](https://opus-codec.org/) — IETF RFC 6716
- [modernc.org/ccgo](https://pkg.go.dev/modernc.org/ccgo/v4) — C to Go transpiler by Jan Mercl
- [modernc.org/libc](https://pkg.go.dev/modernc.org/libc) — Go libc runtime
