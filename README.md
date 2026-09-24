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

| Library | Type | SILK | CELT | Hybrid | CGo |
|---|---|---|---|---|---|
| **go-opus (this)** | Pure Go | Yes | Yes | Yes | No |
| [pion/opus](https://github.com/pion/opus) | Pure Go | Yes | No | No | No |
| [hraban/opus](https://github.com/hraban/opus) | CGo wrapper | Yes | Yes | Yes | Yes |

## Platform Support

Currently built for `darwin/arm64`. To add other platforms, re-run the ccgo transpilation on the target platform.

## License

BSD 3-Clause — same as the original Opus reference implementation (IETF Trust, Skype Limited, Xiph.Org Foundation).

## Credits

- [Opus Codec](https://opus-codec.org/) — IETF RFC 6716
- [modernc.org/ccgo](https://pkg.go.dev/modernc.org/ccgo/v4) — C to Go transpiler by Jan Mercl
- [modernc.org/libc](https://pkg.go.dev/modernc.org/libc) — Go libc runtime
