// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.

package opus

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/skrashevich/go-opus/internal/libc"
)

const demoMaxPacket = 1500

func demoUsage(prog string) {
	fmt.Fprintf(os.Stderr, "Usage: %s [-e] <application> <sampling rate (Hz)> <channels (1/2)> <bits per second> [options] <input> <output>\n", prog)
	fmt.Fprintf(os.Stderr, "       %s -d <sampling rate (Hz)> <channels (1/2)> [options] <input> <output>\n\n", prog)
	fmt.Fprint(os.Stderr, "mode: voip | audio | restricted-lowdelay\n")
	fmt.Fprint(os.Stderr, "options:\n")
	fmt.Fprint(os.Stderr, "-e                   : only runs the encoder (output the bit-stream)\n")
	fmt.Fprint(os.Stderr, "-d                   : only runs the decoder (reads the bit-stream as input)\n")
	fmt.Fprint(os.Stderr, "-cbr                 : enable constant bitrate; default: variable bitrate\n")
	fmt.Fprint(os.Stderr, "-cvbr                : enable constrained variable bitrate; default: unconstrained\n")
	fmt.Fprint(os.Stderr, "-bandwidth <NB|MB|WB|SWB|FB> : audio bandwidth (from narrowband to fullband); default: sampling rate\n")
	fmt.Fprint(os.Stderr, "-framesize <2.5|5|10|20|40|60> : frame size in ms; default: 20\n")
	fmt.Fprint(os.Stderr, "-max_payload <bytes> : maximum payload size in bytes, default: 1024\n")
	fmt.Fprint(os.Stderr, "-complexity <comp>   : complexity, 0 (lowest) ... 10 (highest); default: 10\n")
	fmt.Fprint(os.Stderr, "-inbandfec           : enable SILK inband FEC\n")
	fmt.Fprint(os.Stderr, "-forcemono           : force mono encoding, even for stereo input\n")
	fmt.Fprint(os.Stderr, "-dtx                 : enable SILK DTX\n")
	fmt.Fprint(os.Stderr, "-loss <perc>         : simulate packet loss, in percent (0-100); default: 0\n")
}

type demoConfig struct {
	encodeOnly     bool
	decodeOnly     bool
	application    int32
	sampleRate     int32
	channels       int32
	bitrate        int32
	bandwidth      int32
	frameSize      int32
	maxPayload     int32
	complexity     int32
	useVBR         int32
	cvbr           int32
	useInbandFEC   int32
	forceChannels  int32
	useDTX         int32
	packetLossPerc int32
	input          string
	output         string
}

// RunDemo runs the opus_demo command-line frontend using the codec implementation
// in this package.
func RunDemo(args []string) int {
	if len(args) == 0 {
		args = []string{"opus_demo"}
	}

	cfg, ok := parseDemoArgs(args)
	if !ok {
		return int(EXIT_FAILURE)
	}

	tls := libc.NewTLS()
	defer tls.Close()

	fmt.Fprintf(os.Stderr, "%s\n", demoCString(opus_get_version_string(tls)))
	if cfg.decodeOnly {
		return runDemoDecode(tls, cfg)
	}
	return runDemoEncode(tls, cfg)
}

func parseDemoArgs(args []string) (demoConfig, bool) {
	cfg := demoConfig{
		application:   OPUS_APPLICATION_AUDIO,
		bandwidth:     -1000,
		maxPayload:    1024,
		complexity:    10,
		useVBR:        1,
		forceChannels: -1000,
	}
	if len(args) < 5 {
		demoUsage(args[0])
		return cfg, false
	}

	pos := 1
	if args[pos] == "-e" {
		cfg.encodeOnly = true
		pos++
	} else if args[pos] == "-d" {
		cfg.decodeOnly = true
		pos++
	}

	if !cfg.decodeOnly {
		if len(args)-pos < 6 {
			demoUsage(args[0])
			return cfg, false
		}
		switch args[pos] {
		case "voip":
			cfg.application = OPUS_APPLICATION_VOIP
		case "audio":
			cfg.application = OPUS_APPLICATION_AUDIO
		case "restricted-lowdelay":
			cfg.application = OPUS_APPLICATION_RESTRICTED_LOWDELAY
		default:
			fmt.Fprintf(os.Stderr, "unknown application: %s\n", args[pos])
			demoUsage(args[0])
			return cfg, false
		}
		pos++
	}

	var ok bool
	if cfg.sampleRate, ok = parseInt32(args[pos], "sampling rate"); !ok {
		return cfg, false
	}
	pos++
	if cfg.channels, ok = parseInt32(args[pos], "channels"); !ok {
		return cfg, false
	}
	pos++
	if cfg.channels != 1 && cfg.channels != 2 {
		fmt.Fprintln(os.Stderr, "channels must be 1 or 2")
		return cfg, false
	}
	switch cfg.sampleRate {
	case 8000, 12000, 16000, 24000, 48000:
	default:
		fmt.Fprintln(os.Stderr, "Supported sampling rates are 8000, 12000, 16000, 24000 and 48000.")
		return cfg, false
	}
	cfg.frameSize = cfg.sampleRate / 50

	if !cfg.decodeOnly {
		if cfg.bitrate, ok = parseInt32(args[pos], "bits per second"); !ok {
			return cfg, false
		}
		pos++
	}

	for pos < len(args)-2 {
		opt := strings.ToLower(args[pos])
		switch opt {
		case "-cbr":
			if !checkDemoEncoderOption(cfg.decodeOnly, args[pos]) {
				return cfg, false
			}
			cfg.useVBR = 0
			pos++
		case "-cvbr":
			if !checkDemoEncoderOption(cfg.decodeOnly, args[pos]) {
				return cfg, false
			}
			cfg.cvbr = 1
			pos++
		case "-bandwidth":
			if !checkDemoEncoderOption(cfg.decodeOnly, args[pos]) || !needValue(args, pos) {
				return cfg, false
			}
			switch args[pos+1] {
			case "NB":
				cfg.bandwidth = OPUS_BANDWIDTH_NARROWBAND
			case "MB":
				cfg.bandwidth = OPUS_BANDWIDTH_MEDIUMBAND
			case "WB":
				cfg.bandwidth = OPUS_BANDWIDTH_WIDEBAND
			case "SWB":
				cfg.bandwidth = OPUS_BANDWIDTH_SUPERWIDEBAND
			case "FB":
				cfg.bandwidth = OPUS_BANDWIDTH_FULLBAND
			default:
				fmt.Fprintf(os.Stderr, "Unknown bandwidth %s. Supported are NB, MB, WB, SWB, FB.\n", args[pos+1])
				return cfg, false
			}
			pos += 2
		case "-framesize":
			if !checkDemoEncoderOption(cfg.decodeOnly, args[pos]) || !needValue(args, pos) {
				return cfg, false
			}
			switch args[pos+1] {
			case "2.5":
				cfg.frameSize = cfg.sampleRate / 400
			case "5":
				cfg.frameSize = cfg.sampleRate / 200
			case "10":
				cfg.frameSize = cfg.sampleRate / 100
			case "20":
				cfg.frameSize = cfg.sampleRate / 50
			case "40":
				cfg.frameSize = cfg.sampleRate / 25
			case "60":
				cfg.frameSize = 3 * cfg.sampleRate / 50
			default:
				fmt.Fprintf(os.Stderr, "Unsupported frame size: %s ms. Supported are 2.5, 5, 10, 20, 40, 60.\n", args[pos+1])
				return cfg, false
			}
			pos += 2
		case "-max_payload":
			if !checkDemoEncoderOption(cfg.decodeOnly, args[pos]) || !needValue(args, pos) {
				return cfg, false
			}
			if cfg.maxPayload, ok = parseInt32(args[pos+1], "max_payload"); !ok {
				return cfg, false
			}
			pos += 2
		case "-complexity":
			if !checkDemoEncoderOption(cfg.decodeOnly, args[pos]) || !needValue(args, pos) {
				return cfg, false
			}
			if cfg.complexity, ok = parseInt32(args[pos+1], "complexity"); !ok {
				return cfg, false
			}
			pos += 2
		case "-inbandfec":
			cfg.useInbandFEC = 1
			pos++
		case "-forcemono":
			if !checkDemoEncoderOption(cfg.decodeOnly, args[pos]) {
				return cfg, false
			}
			cfg.forceChannels = 1
			pos++
		case "-dtx":
			if !checkDemoEncoderOption(cfg.decodeOnly, args[pos]) {
				return cfg, false
			}
			cfg.useDTX = 1
			pos++
		case "-loss":
			if !checkDemoDecoderOption(cfg.encodeOnly, args[pos]) || !needValue(args, pos) {
				return cfg, false
			}
			if cfg.packetLossPerc, ok = parseInt32(args[pos+1], "loss"); !ok {
				return cfg, false
			}
			pos += 2
		default:
			fmt.Fprintf(os.Stderr, "Error: unrecognized setting: %s\n\n", args[pos])
			demoUsage(args[0])
			return cfg, false
		}
	}

	if cfg.maxPayload < 0 || cfg.maxPayload > demoMaxPacket {
		fmt.Fprintf(os.Stderr, "max_payload_bytes must be between 0 and %d\n", demoMaxPacket)
		return cfg, false
	}
	cfg.input = args[len(args)-2]
	cfg.output = args[len(args)-1]
	return cfg, true
}

func runDemoEncode(tls *libc.TLS, cfg demoConfig) int {
	input, err := os.ReadFile(cfg.input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file %s\n", cfg.input)
		return int(EXIT_FAILURE)
	}

	var encErr int32
	enc := opus_encoder_create(tls, cfg.sampleRate, cfg.channels, cfg.application, uintptr(unsafe.Pointer(&encErr)))
	if encErr != OPUS_OK {
		fmt.Fprintf(os.Stderr, "Cannot create encoder: %s\n", demoCString(opus_strerror(tls, encErr)))
		return int(EXIT_FAILURE)
	}
	defer opus_encoder_destroy(tls, enc)

	ctl1(tls, enc, OPUS_SET_BITRATE_REQUEST, cfg.bitrate)
	ctl1(tls, enc, OPUS_SET_BANDWIDTH_REQUEST, cfg.bandwidth)
	ctl1(tls, enc, OPUS_SET_VBR_REQUEST, cfg.useVBR)
	ctl1(tls, enc, OPUS_SET_VBR_CONSTRAINT_REQUEST, cfg.cvbr)
	ctl1(tls, enc, OPUS_SET_COMPLEXITY_REQUEST, cfg.complexity)
	ctl1(tls, enc, OPUS_SET_INBAND_FEC_REQUEST, cfg.useInbandFEC)
	ctl1(tls, enc, OPUS_SET_FORCE_CHANNELS_REQUEST, cfg.forceChannels)
	ctl1(tls, enc, OPUS_SET_DTX_REQUEST, cfg.useDTX)
	ctl1(tls, enc, OPUS_SET_PACKET_LOSS_PERC_REQUEST, cfg.packetLossPerc)

	var skip int32
	encCtlPtr(tls, enc, OPUS_GET_LOOKAHEAD_REQUEST, uintptr(unsafe.Pointer(&skip)))

	fmt.Fprintf(os.Stderr, "Encoding %d Hz input at %.3f kb/s in %s mode with %d-sample frames.\n",
		cfg.sampleRate, float64(cfg.bitrate)*0.001, bandwidthName(cfg.bandwidth), cfg.frameSize)

	packetized, decoded, stats, ok := encodeFrames(tls, cfg, enc, input, skip)
	if !ok {
		return int(EXIT_FAILURE)
	}

	var out []byte
	if cfg.encodeOnly {
		out = packetized
	} else {
		out = decoded
	}
	if err := os.WriteFile(cfg.output, out, 0666); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing.")
		return int(EXIT_FAILURE)
	}
	printDemoStats(cfg.sampleRate, cfg.frameSize, stats)
	return int(EXIT_SUCCESS)
}

func runDemoDecode(tls *libc.TLS, cfg demoConfig) int {
	input, err := os.ReadFile(cfg.input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file %s\n", cfg.input)
		return int(EXIT_FAILURE)
	}

	var decErr int32
	dec := opus_decoder_create(tls, cfg.sampleRate, cfg.channels, uintptr(unsafe.Pointer(&decErr)))
	if decErr != OPUS_OK {
		fmt.Fprintf(os.Stderr, "Cannot create decoder: %s\n", demoCString(opus_strerror(tls, decErr)))
		return int(EXIT_FAILURE)
	}
	defer opus_decoder_destroy(tls, dec)

	fmt.Fprintf(os.Stderr, "Decoding with %d Hz output (%d channels)\n", cfg.sampleRate, cfg.channels)

	decoded, stats, ok := decodePackets(tls, cfg, dec, input, 0)
	if !ok {
		return int(EXIT_FAILURE)
	}
	if err := os.WriteFile(cfg.output, decoded, 0666); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing.")
		return int(EXIT_FAILURE)
	}
	printDemoStats(cfg.sampleRate, cfg.frameSize, stats)
	return int(EXIT_SUCCESS)
}

type demoStats struct {
	frames     int
	bits       float64
	bits2      float64
	bitsMax    float64
	active     int
	activeBits float64
}

func encodeFrames(tls *libc.TLS, cfg demoConfig, enc uintptr, input []byte, skip int32) ([]byte, []byte, demoStats, bool) {
	var dec uintptr
	if !cfg.encodeOnly {
		var decErr int32
		dec = opus_decoder_create(tls, cfg.sampleRate, cfg.channels, uintptr(unsafe.Pointer(&decErr)))
		if decErr != OPUS_OK {
			fmt.Fprintf(os.Stderr, "Cannot create decoder: %s\n", demoCString(opus_strerror(tls, decErr)))
			return nil, nil, demoStats{}, false
		}
		defer opus_decoder_destroy(tls, dec)
	}

	frameSamples := int(cfg.frameSize * cfg.channels)
	frameBytes := frameSamples * 2
	pcm := make([]opus_int16, frameSamples)
	packet := make([]byte, cfg.maxPayload)
	decoded := make([]opus_int16, int(cfg.sampleRate/25*cfg.channels*6))
	var packetized []byte
	var decodedBytes []byte
	stats := demoStats{}

	for offset := 0; offset < len(input); offset += frameBytes {
		chunk := input[offset:min(offset+frameBytes, len(input))]
		clear(pcm)
		for i := 0; i+1 < len(chunk); i += 2 {
			pcm[i/2] = opus_int16(int16(binary.LittleEndian.Uint16(chunk[i:])))
		}

		n := opus_encode(tls, enc, demoPtrOf(pcm), cfg.frameSize, demoPtrOf(packet), cfg.maxPayload)
		if n < 0 {
			fmt.Fprintf(os.Stderr, "opus_encode() returned %d\n", n)
			return nil, nil, stats, false
		}
		var finalRange opus_uint32
		encCtlPtr(tls, enc, OPUS_GET_FINAL_RANGE_REQUEST, uintptr(unsafe.Pointer(&finalRange)))
		packetized = appendPacket(packetized, uint32(n), uint32(finalRange), packet[:n])
		updateDemoStats(&stats, int(n), pcm)

		if dec != 0 {
			lost := cfg.packetLossPerc > 0 && rand.Int31n(100) < cfg.packetLossPerc
			var dataPtr uintptr
			if !lost && n > 0 {
				dataPtr = demoPtrOf(packet[:n])
			}
			outputSamples := opus_decode(tls, dec, dataPtr, n, demoPtrOf(decoded), int32(len(decoded))/cfg.channels, 0)
			if outputSamples < 0 {
				fmt.Fprintf(os.Stderr, "error decoding frame: %s\n", demoCString(opus_strerror(tls, outputSamples)))
				return nil, nil, stats, false
			}
			var decFinalRange opus_uint32
			decCtlPtr(tls, dec, OPUS_GET_FINAL_RANGE_REQUEST, uintptr(unsafe.Pointer(&decFinalRange)))
			if !lost && finalRange != decFinalRange {
				fmt.Fprintf(os.Stderr, "Error: Range coder state mismatch between encoder and decoder in frame %d: 0x%8x vs 0x%8x\n", stats.frames-1, finalRange, decFinalRange)
				return nil, nil, stats, false
			}
			start := int(skip) * int(cfg.channels)
			if start > int(outputSamples*cfg.channels) {
				skip -= outputSamples
				continue
			}
			skip = 0
			decodedBytes = appendPCM(decodedBytes, decoded[start:int(outputSamples*cfg.channels)])
		}
	}

	return packetized, decodedBytes, stats, true
}

func decodePackets(tls *libc.TLS, cfg demoConfig, dec uintptr, input []byte, skip int32) ([]byte, demoStats, bool) {
	decoded := make([]opus_int16, int(cfg.sampleRate/25*cfg.channels*6))
	var decodedBytes []byte
	stats := demoStats{}
	for offset := 0; offset < len(input); {
		if offset+8 > len(input) {
			fmt.Fprintln(os.Stderr, "Invalid payload length: truncated packet header")
			return nil, stats, false
		}
		n := int32(binary.BigEndian.Uint32(input[offset:]))
		finalRange := opus_uint32(binary.BigEndian.Uint32(input[offset+4:]))
		offset += 8
		if n < 0 || n > cfg.maxPayload || offset+int(n) > len(input) {
			fmt.Fprintf(os.Stderr, "Invalid payload length: %d\n", n)
			return nil, stats, false
		}
		packet := input[offset : offset+int(n)]
		offset += int(n)

		lost := cfg.packetLossPerc > 0 && rand.Int31n(100) < cfg.packetLossPerc
		var dataPtr uintptr
		if !lost && n > 0 {
			dataPtr = demoPtrOf(packet)
		}
		outputSamples := opus_decode(tls, dec, dataPtr, n, demoPtrOf(decoded), int32(len(decoded))/cfg.channels, 0)
		if outputSamples < 0 {
			fmt.Fprintf(os.Stderr, "error decoding frame: %s\n", demoCString(opus_strerror(tls, outputSamples)))
			return nil, stats, false
		}
		var decFinalRange opus_uint32
		decCtlPtr(tls, dec, OPUS_GET_FINAL_RANGE_REQUEST, uintptr(unsafe.Pointer(&decFinalRange)))
		if !lost && finalRange != 0 && finalRange != decFinalRange {
			fmt.Fprintf(os.Stderr, "Error: Range coder state mismatch between encoder and decoder in frame %d: 0x%8x vs 0x%8x\n", stats.frames, finalRange, decFinalRange)
			return nil, stats, false
		}
		updateDemoStats(&stats, int(n), nil)

		start := int(skip) * int(cfg.channels)
		if start > int(outputSamples*cfg.channels) {
			skip -= outputSamples
			continue
		}
		skip = 0
		decodedBytes = appendPCM(decodedBytes, decoded[start:int(outputSamples*cfg.channels)])
	}
	return decodedBytes, stats, true
}

func appendPacket(dst []byte, n uint32, finalRange uint32, data []byte) []byte {
	var header [8]byte
	binary.BigEndian.PutUint32(header[:4], n)
	binary.BigEndian.PutUint32(header[4:], finalRange)
	dst = append(dst, header[:]...)
	return append(dst, data...)
}

func appendPCM(dst []byte, pcm []opus_int16) []byte {
	var sample [2]byte
	for _, v := range pcm {
		binary.LittleEndian.PutUint16(sample[:], uint16(v))
		dst = append(dst, sample[:]...)
	}
	return dst
}

func updateDemoStats(stats *demoStats, n int, pcm []opus_int16) {
	bits := float64(n * 8)
	stats.frames++
	stats.bits += bits
	stats.bits2 += bits * bits
	if bits > stats.bitsMax {
		stats.bitsMax = bits
	}

	var energy float64
	for _, sample := range pcm {
		energy += float64(sample) * float64(sample)
	}
	if len(pcm) > 0 && energy/float64(len(pcm)) > 100000 {
		stats.active++
		stats.activeBits += bits
	}
}

func printDemoStats(sampleRate, frameSize int32, stats demoStats) {
	if stats.frames == 0 {
		return
	}
	denom := float64(frameSize * int32(stats.frames))
	fmt.Fprintf(os.Stderr, "average bitrate:             %7.3f kb/s\n", 0.001*stats.bits*float64(sampleRate)/denom)
	fmt.Fprintf(os.Stderr, "maximum bitrate:             %7.3f bkp/s\n", 0.001*stats.bitsMax*float64(sampleRate)/float64(frameSize))
	if stats.active > 0 {
		fmt.Fprintf(os.Stderr, "active bitrate:              %7.3f kb/s\n", 0.001*stats.activeBits*float64(sampleRate)/float64(frameSize*int32(stats.active)))
	}
}

func ctl1(tls *libc.TLS, enc uintptr, request int32, value int32) {
	bp := tls.Alloc(16)
	defer tls.Free(16)
	opus_encoder_ctl(tls, enc, request, libc.VaList(bp, value))
}

func encCtlPtr(tls *libc.TLS, state uintptr, request int32, ptr uintptr) {
	bp := tls.Alloc(16)
	defer tls.Free(16)
	opus_encoder_ctl(tls, state, request, libc.VaList(bp, ptr))
}

func decCtlPtr(tls *libc.TLS, state uintptr, request int32, ptr uintptr) {
	bp := tls.Alloc(16)
	defer tls.Free(16)
	opus_decoder_ctl(tls, state, request, libc.VaList(bp, ptr))
}

func parseInt32(s, name string) (int32, bool) {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid %s: %s\n", name, s)
		return 0, false
	}
	return int32(v), true
}

func needValue(args []string, pos int) bool {
	if pos+1 >= len(args)-2 {
		fmt.Fprintf(os.Stderr, "missing value for %s\n", args[pos])
		return false
	}
	return true
}

func checkDemoEncoderOption(decodeOnly bool, opt string) bool {
	if decodeOnly {
		fmt.Fprintf(os.Stderr, "option %s is only for encoding\n", opt)
		return false
	}
	return true
}

func checkDemoDecoderOption(encodeOnly bool, opt string) bool {
	if encodeOnly {
		fmt.Fprintf(os.Stderr, "option %s is only for decoding\n", opt)
		return false
	}
	return true
}

func bandwidthName(bandwidth int32) string {
	switch bandwidth {
	case OPUS_BANDWIDTH_NARROWBAND:
		return "narrowband"
	case OPUS_BANDWIDTH_MEDIUMBAND:
		return "mediumband"
	case OPUS_BANDWIDTH_WIDEBAND:
		return "wideband"
	case OPUS_BANDWIDTH_SUPERWIDEBAND:
		return "superwideband"
	case OPUS_BANDWIDTH_FULLBAND:
		return "fullband"
	case -1000:
		return "auto"
	default:
		return "unknown"
	}
}

func demoCString(p uintptr) string {
	if p == 0 {
		return ""
	}
	var b []byte
	for i := uintptr(0); ; i++ {
		c := *(*byte)(unsafe.Pointer(p + i))
		if c == 0 {
			return string(b)
		}
		b = append(b, c)
	}
}

func demoPtrOf[T any](v []T) uintptr {
	if len(v) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&v[0]))
}
