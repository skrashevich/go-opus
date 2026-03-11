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

Benchmarked on Apple M1 Pro, encoding/decoding a 3:41 stereo 48kHz file at 128 kbps (median of 3 runs, user CPU time):

| Implementation | Encode | Decode | Notes |
|---|---|---|---|
| C original (RFC 6716) | 1.51s (1.0x) | 0.60s (1.0x) | Baseline |
| **Go transpiled (this)** | **2.64s (1.7x)** | **0.96s (1.6x)** | Pure Go, no CGo |
| FFmpeg libopus | 1.84s (1.2x) | 0.54s (0.9x) | SIMD-optimized |

~1.7x slower than C — excellent for automatic transpilation. Encode runs at 84x realtime, decode at 230x realtime.

### Detailed benchmarks (3 runs, user CPU time)

| Implementation | Run 1 | Run 2 | Run 3 | Median | vs C |
|---|---|---|---|---|---|
| **Encode** | | | | | |
| C original | 1.50s | 1.53s | 1.51s | 1.51s | 1.0x |
| Go transpiled | 2.67s | 2.64s | 2.64s | 2.64s | 1.7x |
| FFmpeg libopus | 1.84s | 1.84s | 1.87s | 1.84s | 1.2x |
| **Decode** | | | | | |
| C original | 0.60s | 0.61s | 0.60s | 0.60s | 1.0x |
| Go transpiled | 0.96s | 0.96s | 0.96s | 0.96s | 1.6x |
| FFmpeg libopus | 0.54s | 0.54s | 0.54s | 0.54s | 0.9x |

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
├── lib_darwin_arm64.go          # Opus library (package opus)
├── cmd/
│   ├── opus_demo/               # Encoder/decoder CLI
│   │   └── main_darwin_arm64.go
│   └── opus_compare/            # Audio comparison CLI
│       └── main_darwin_arm64.go
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
