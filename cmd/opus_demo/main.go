// Command opus_demo is a Go port of src/opus_demo.c from libopus. The codec
// itself is the github.com/skrashevich/go-opus package; everything else (option
// parsing, file I/O, statistics) is plain Go that reproduces the output of the
// C program byte for byte.
//
// Original C source:
//
//	Copyright (c) 2007-2008 CSIRO
//	Copyright (c) 2007-2009 Xiph.Org Foundation
//	Written by Jean-Marc Valin
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// - Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
// - Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES ARE DISCLAIMED. SEE THE ORIGINAL FILE
// FOR THE FULL LICENSE TEXT.
package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"runtime/pprof"
	"unsafe"

	opus "github.com/skrashevich/go-opus"
)

const (
	exitSuccess = 0
	exitFailure = 1

	maxPacket = 1500

	applicationVoIP               = 2048
	applicationAudio              = 2049
	applicationRestrictedLowdelay = 2051

	bandwidthNarrowband    = 1101
	bandwidthMediumband    = 1102
	bandwidthWideband      = 1103
	bandwidthSuperwideband = 1104
	bandwidthFullband      = 1105

	setBitrateRequest        = 4002
	setVBRRequest            = 4006
	setBandwidthRequest      = 4008
	setComplexityRequest     = 4010
	setInbandFECRequest      = 4012
	setPacketLossPercRequest = 4014
	setDTXRequest            = 4016
	setVBRConstraintRequest  = 4020
	setForceChannelsRequest  = 4022
	getLookaheadRequest      = 4027
	getFinalRangeRequest     = 4031
	setForceModeRequest      = 11002

	modeSilkOnly = 1000
	modeCeltOnly = 1002
)

// stopProfile flushes the CPU profile requested via $CPUPROFILE, if any.
var stopProfile = func() {}

func exit(code int) {
	stopProfile()
	os.Exit(code)
}

func main() {
	if p := os.Getenv("CPUPROFILE"); p != "" {
		f, _ := os.Create(p)
		pprof.StartCPUProfile(f)
		stopProfile = func() {
			pprof.StopCPUProfile()
			f.Close()
		}
	}
	exit(run(os.Args))
}

// nextRand is the state of crand.
var nextRand uint64 = 1

// crand reproduces rand(3) of modernc.org/libc on non-Linux targets (libc.go),
// which the reference binary used on darwin: a 64-bit LCG
// next = next*1103515245 + 12345 with implicit seed 1, returning
// (next >> 32) % 0x7fffffff. (modernc's Linux build uses musl's rand instead,
// seed = seed*6364136223846793005 + 1 returning seed >> 33, starting at 0.)
func crand() int32 {
	nextRand = nextRand*1103515245 + 12345
	return int32(uint32(nextRand / (math.MaxUint32 + 1) % math.MaxInt32))
}

// atol reproduces atol/atoi of modernc.org/libc: optional leading whitespace,
// one optional sign, then decimal digits accumulated in a uint64 until the
// first non-digit or a detected overflow. Callers truncate to int32 as C does.
func atol(s string) int64 {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r' || s[i] == '\v' || s[i] == '\f') {
		i++
	}
	neg := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		neg = s[i] == '-'
		i++
	}
	var n uint64
	for ; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		n0 := n
		n = 10*n + uint64(s[i]-'0')
		if n < n0 {
			n = n0
			break
		}
	}
	if neg {
		return int64(-n)
	}
	return int64(n)
}

func atoi(s string) int32 { return int32(atol(s)) }

// strcaseeq reports whether a and b are equal ignoring ASCII case, like
// strcasecmp(a, b) == 0.
func strcaseeq(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		x, y := a[i], b[i]
		if x >= 'a' && x <= 'z' {
			x -= 'a' - 'A'
		}
		if y >= 'a' && y <= 'z' {
			y -= 'a' - 'A'
		}
		if x != y {
			return false
		}
	}
	return true
}

// cptr converts the address of C memory (outside the Go heap, from
// opus.Malloc or returned by the codec) to an unsafe.Pointer.
func cptr(p uintptr) unsafe.Pointer { return *(*unsafe.Pointer)(unsafe.Pointer(&p)) }

// cmalloc returns n zeroed bytes of C memory.
func cmalloc(n int) uintptr {
	if n < 1 {
		n = 1
	}
	p := opus.Malloc(n)
	if p == 0 {
		panic("opus_demo: out of memory")
	}
	clear(unsafe.Slice((*byte)(cptr(p)), n))
	return p
}

// goString converts a NUL-terminated C string to a Go string.
func goString(p uintptr) string {
	if p == 0 {
		return ""
	}
	n := 0
	for *(*byte)(cptr(p + uintptr(n))) != 0 {
		n++
	}
	return string(unsafe.Slice((*byte)(cptr(p)), n))
}

func printUsage(argv0 string) {
	fmt.Fprintf(os.Stderr, "Usage: %s [-e] <application> <sampling rate (Hz)> <channels (1/2)> <bits per second>  [options] <input> <output>\n", argv0)
	fmt.Fprintf(os.Stderr, "       %s -d <sampling rate (Hz)> <channels (1/2)> [options] <input> <output>\n\n", argv0)
	fmt.Fprint(os.Stderr, "mode: voip | audio | restricted-lowdelay\n")
	fmt.Fprint(os.Stderr, "options:\n")
	fmt.Fprint(os.Stderr, "-e                   : only runs the encoder (output the bit-stream)\n")
	fmt.Fprint(os.Stderr, "-d                   : only runs the decoder (reads the bit-stream as input)\n")
	fmt.Fprint(os.Stderr, "-cbr                 : enable constant bitrate; default: variable bitrate\n")
	fmt.Fprint(os.Stderr, "-cvbr                : enable constrained variable bitrate; default: unconstrained\n")
	fmt.Fprint(os.Stderr, "-bandwidth <NB|MB|WB|SWB|FB> : audio bandwidth (from narrowband to fullband); default: sampling rate\n")
	fmt.Fprint(os.Stderr, "-framesize <2.5|5|10|20|40|60> : frame size in ms; default: 20 \n")
	fmt.Fprint(os.Stderr, "-max_payload <bytes> : maximum payload size in bytes, default: 1024\n")
	fmt.Fprint(os.Stderr, "-complexity <comp>   : complexity, 0 (lowest) ... 10 (highest); default: 10\n")
	fmt.Fprint(os.Stderr, "-inbandfec           : enable SILK inband FEC\n")
	fmt.Fprint(os.Stderr, "-forcemono           : force mono encoding, even for stereo input\n")
	fmt.Fprint(os.Stderr, "-dtx                 : enable SILK DTX\n")
	fmt.Fprint(os.Stderr, "-loss <perc>         : simulate packet loss, in percent (0-100); default: 0\n")
}

func intToChar(i uint32, ch []byte) {
	ch[0] = byte(i >> 24)
	ch[1] = byte(i >> 16)
	ch[2] = byte(i >> 8)
	ch[3] = byte(i)
}

func charToInt(ch []byte) uint32 {
	return uint32(ch[0])<<24 | uint32(ch[1])<<16 | uint32(ch[2])<<8 | uint32(ch[3])
}

func checkDecoderOption(encodeOnly bool, opt string) {
	if encodeOnly {
		fmt.Fprintf(os.Stderr, "option %s is only for decoding\n", opt)
		exit(exitFailure)
	}
}

func checkEncoderOption(decodeOnly bool, opt string) {
	if decodeOnly {
		fmt.Fprintf(os.Stderr, "option %s is only for encoding\n", opt)
		exit(exitFailure)
	}
}

// testMode is one entry of the mode-switching test lists.
type testMode struct {
	mode, bandwidth, frameSize, channels int32
}

var silk8Test = []testMode{
	{modeSilkOnly, bandwidthNarrowband, 960 * 3, 1},
	{modeSilkOnly, bandwidthNarrowband, 960 * 2, 1},
	{modeSilkOnly, bandwidthNarrowband, 960, 1},
	{modeSilkOnly, bandwidthNarrowband, 480, 1},
	{modeSilkOnly, bandwidthNarrowband, 960 * 3, 2},
	{modeSilkOnly, bandwidthNarrowband, 960 * 2, 2},
	{modeSilkOnly, bandwidthNarrowband, 960, 2},
	{modeSilkOnly, bandwidthNarrowband, 480, 2},
}

var silk12Test = []testMode{
	{modeSilkOnly, bandwidthMediumband, 960 * 3, 1},
	{modeSilkOnly, bandwidthMediumband, 960 * 2, 1},
	{modeSilkOnly, bandwidthMediumband, 960, 1},
	{modeSilkOnly, bandwidthMediumband, 480, 1},
	{modeSilkOnly, bandwidthMediumband, 960 * 3, 2},
	{modeSilkOnly, bandwidthMediumband, 960 * 2, 2},
	{modeSilkOnly, bandwidthMediumband, 960, 2},
	{modeSilkOnly, bandwidthMediumband, 480, 2},
}

var silk16Test = []testMode{
	{modeSilkOnly, bandwidthWideband, 960 * 3, 1},
	{modeSilkOnly, bandwidthWideband, 960 * 2, 1},
	{modeSilkOnly, bandwidthWideband, 960, 1},
	{modeSilkOnly, bandwidthWideband, 480, 1},
	{modeSilkOnly, bandwidthWideband, 960 * 3, 2},
	{modeSilkOnly, bandwidthWideband, 960 * 2, 2},
	{modeSilkOnly, bandwidthWideband, 960, 2},
	{modeSilkOnly, bandwidthWideband, 480, 2},
}

var hybrid24Test = []testMode{
	{modeSilkOnly, bandwidthSuperwideband, 960, 1},
	{modeSilkOnly, bandwidthSuperwideband, 480, 1},
	{modeSilkOnly, bandwidthSuperwideband, 960, 2},
	{modeSilkOnly, bandwidthSuperwideband, 480, 2},
}

var hybrid48Test = []testMode{
	{modeSilkOnly, bandwidthFullband, 960, 1},
	{modeSilkOnly, bandwidthFullband, 480, 1},
	{modeSilkOnly, bandwidthFullband, 960, 2},
	{modeSilkOnly, bandwidthFullband, 480, 2},
}

var celtTest = func() []testMode {
	var l []testMode
	for _, ch := range []int32{1, 2} {
		for _, fs := range []int32{960, 480, 240, 120} {
			for _, bw := range []int32{bandwidthFullband, bandwidthSuperwideband, bandwidthWideband, bandwidthNarrowband} {
				l = append(l, testMode{modeCeltOnly, bw, fs, ch})
			}
		}
	}
	return l
}()

var celtHQTest = []testMode{
	{modeCeltOnly, bandwidthFullband, 960, 2},
	{modeCeltOnly, bandwidthFullband, 480, 2},
	{modeCeltOnly, bandwidthFullband, 240, 2},
	{modeCeltOnly, bandwidthFullband, 120, 2},
}

// readFull mimics fread(buf, 1, len(buf), f): it returns the number of bytes
// read, stopping early only at end of file or on a read error.
func readFull(r io.Reader, buf []byte) int {
	n, _ := io.ReadFull(r, buf)
	return n
}

func run(argv []string) int {
	argc := int32(len(argv))
	var (
		bitrateBps              int32
		cvbr                    int32
		count, countAct         int32
		stop                    bool
		application             int32 = applicationAudio
		bits, bitsMax, bitsAct  float64
		bits2, nrg              float64
		bandwidth               int32
		lost                    bool
		lostPrev                = true
		toggle                  int32
		lens                    [2]int32
		encodeOnly, decodeOnly  bool
		maxFrameSize            int32 = 960 * 6
		currRead                int32
		sweepBps                int32
		randomFramesize         bool
		newsize                 int32
		delayedCelt             bool
		sweepMax, sweepMin      int32
		randomFEC               bool
		modeList                []testMode
		currMode, currModeCount int32
		modeSwitchTime          int32 = 48000
		useInbandFEC            int32
		enc, dec                uintptr
		outputSamples           int32
	)

	if argc < 5 {
		printUsage(argv[0])
		return exitFailure
	}

	tls := opus.NewTLS()
	defer tls.Close()

	fmt.Fprintf(os.Stderr, "%s\n", goString(opus.GetVersionString(tls)))

	args := int32(1)
	if argv[args] == "-e" {
		encodeOnly = true
		args++
	} else if argv[args] == "-d" {
		decodeOnly = true
		args++
	}
	if !decodeOnly && argc < 7 {
		printUsage(argv[0])
		return exitFailure
	}

	if !decodeOnly {
		switch argv[args] {
		case "voip":
			application = applicationVoIP
		case "restricted-lowdelay":
			application = applicationRestrictedLowdelay
		case "audio":
		default:
			fmt.Fprintf(os.Stderr, "unknown application: %s\n", argv[args])
			printUsage(argv[0])
			return exitFailure
		}
		args++
	}
	samplingRate := int32(atol(argv[args]))
	args++
	channels := atoi(argv[args])
	args++
	if !decodeOnly {
		bitrateBps = int32(atol(argv[args]))
		args++
	}

	if samplingRate != 8000 && samplingRate != 12000 && samplingRate != 16000 &&
		samplingRate != 24000 && samplingRate != 48000 {
		fmt.Fprint(os.Stderr, "Supported sampling rates are 8000, 12000, 16000, 24000 and 48000.\n")
		return exitFailure
	}
	frameSize := samplingRate / 50

	/* defaults: */
	useVBR := int32(1)
	bandwidth = -1000
	maxPayloadBytes := int32(maxPacket)
	complexity := int32(10)
	forcechannels := int32(-1000)
	useDTX := int32(0)
	packetLossPerc := int32(0)

	for args < argc-2 {
		/* process command line options */
		arg := argv[args]
		switch {
		case strcaseeq(arg, "-cbr"):
			checkEncoderOption(decodeOnly, "-cbr")
			useVBR = 0
			args++
		case strcaseeq(arg, "-bandwidth"):
			checkEncoderOption(decodeOnly, "-bandwidth")
			switch v := argv[args+1]; v {
			case "NB":
				bandwidth = bandwidthNarrowband
			case "MB":
				bandwidth = bandwidthMediumband
			case "WB":
				bandwidth = bandwidthWideband
			case "SWB":
				bandwidth = bandwidthSuperwideband
			case "FB":
				bandwidth = bandwidthFullband
			default:
				fmt.Fprintf(os.Stderr, "Unknown bandwidth %s. Supported are NB, MB, WB, SWB, FB.\n", v)
				return exitFailure
			}
			args += 2
		case strcaseeq(arg, "-framesize"):
			checkEncoderOption(decodeOnly, "-framesize")
			switch v := argv[args+1]; v {
			case "2.5":
				frameSize = samplingRate / 400
			case "5":
				frameSize = samplingRate / 200
			case "10":
				frameSize = samplingRate / 100
			case "20":
				frameSize = samplingRate / 50
			case "40":
				frameSize = samplingRate / 25
			case "60":
				frameSize = 3 * samplingRate / 50
			default:
				fmt.Fprintf(os.Stderr, "Unsupported frame size: %s ms. Supported are 2.5, 5, 10, 20, 40, 60.\n", v)
				return exitFailure
			}
			args += 2
		case strcaseeq(arg, "-max_payload"):
			checkEncoderOption(decodeOnly, "-max_payload")
			maxPayloadBytes = atoi(argv[args+1])
			args += 2
		case strcaseeq(arg, "-complexity"):
			checkEncoderOption(decodeOnly, "-complexity")
			complexity = atoi(argv[args+1])
			args += 2
		case strcaseeq(arg, "-inbandfec"):
			useInbandFEC = 1
			args++
		case strcaseeq(arg, "-forcemono"):
			checkEncoderOption(decodeOnly, "-forcemono")
			forcechannels = 1
			args++
		case strcaseeq(arg, "-cvbr"):
			checkEncoderOption(decodeOnly, "-cvbr")
			cvbr = 1
			args++
		case strcaseeq(arg, "-dtx"):
			checkEncoderOption(decodeOnly, "-dtx")
			useDTX = 1
			args++
		case strcaseeq(arg, "-loss"):
			checkDecoderOption(encodeOnly, "-loss")
			packetLossPerc = atoi(argv[args+1])
			args += 2
		case strcaseeq(arg, "-sweep"):
			checkEncoderOption(decodeOnly, "-sweep")
			sweepBps = atoi(argv[args+1])
			args += 2
		case strcaseeq(arg, "-random_framesize"):
			checkEncoderOption(decodeOnly, "-random_framesize")
			randomFramesize = true
			args++
		case strcaseeq(arg, "-sweep_max"):
			checkEncoderOption(decodeOnly, "-sweep_max")
			sweepMax = atoi(argv[args+1])
			args += 2
		case strcaseeq(arg, "-random_fec"):
			checkEncoderOption(decodeOnly, "-random_fec")
			randomFEC = true
			args++
		case strcaseeq(arg, "-silk8k_test"):
			checkEncoderOption(decodeOnly, "-silk8k_test")
			modeList = silk8Test
			args++
		case strcaseeq(arg, "-silk12k_test"):
			checkEncoderOption(decodeOnly, "-silk12k_test")
			modeList = silk12Test
			args++
		case strcaseeq(arg, "-silk16k_test"):
			checkEncoderOption(decodeOnly, "-silk16k_test")
			modeList = silk16Test
			args++
		case strcaseeq(arg, "-hybrid24k_test"):
			checkEncoderOption(decodeOnly, "-hybrid24k_test")
			modeList = hybrid24Test
			args++
		case strcaseeq(arg, "-hybrid48k_test"):
			checkEncoderOption(decodeOnly, "-hybrid48k_test")
			modeList = hybrid48Test
			args++
		case strcaseeq(arg, "-celt_test"):
			checkEncoderOption(decodeOnly, "-celt_test")
			modeList = celtTest
			args++
		case strcaseeq(arg, "-celt_hq_test"):
			checkEncoderOption(decodeOnly, "-celt_hq_test")
			modeList = celtHQTest
			args++
		default:
			fmt.Printf("Error: unrecognized setting: %s\n\n", arg)
			printUsage(argv[0])
			return exitFailure
		}
	}
	nbModesInList := int32(len(modeList))

	if sweepMax != 0 {
		sweepMin = bitrateBps
	}

	if maxPayloadBytes < 0 || maxPayloadBytes > maxPacket {
		fmt.Fprintf(os.Stderr, "max_payload_bytes must be between 0 and %d\n", maxPacket)
		return exitFailure
	}

	inFile := argv[argc-2]
	finFile, err := os.Open(inFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file %s\n", inFile)
		return exitFailure
	}
	defer finFile.Close()
	if modeList != nil {
		end, _ := finFile.Seek(0, io.SeekEnd)
		size := int32(end)
		fmt.Fprintf(os.Stderr, "File size is %d bytes\n", size)
		finFile.Seek(0, io.SeekStart)
		modeSwitchTime = int32(uint64(size) / 2 / uint64(channels) / uint64(nbModesInList))
		fmt.Fprintf(os.Stderr, "Switching mode every %d samples\n", modeSwitchTime)
	}
	fin := bufio.NewReader(finFile)

	outFile := argv[argc-1]
	foutFile, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open output file %s\n", outFile)
		return exitFailure
	}
	defer foutFile.Close()
	fout := bufio.NewWriter(foutFile)
	defer fout.Flush()
	// fwrite mimics fwrite(buf, 1, len(buf), fout) == len(buf).
	fwrite := func(b []byte) bool {
		n, _ := fout.Write(b)
		return n == len(b)
	}

	// C-side scalars handed to the codec by address: the error code, the
	// encoder lookahead ("skip"), the encoder final ranges, the decoder final
	// range, and one va_list slot for ctl arguments.
	scratch := cmalloc(32)
	defer opus.Free(scratch)
	errPtr, skipPtr, encRangePtr, decRangePtr, va := scratch, scratch+4, scratch+8, scratch+16, scratch+24
	errCode := (*int32)(cptr(errPtr))
	skip := (*int32)(cptr(skipPtr))
	encFinalRange := (*[2]uint32)(cptr(encRangePtr))
	decFinalRange := (*uint32)(cptr(decRangePtr))
	encCtl := func(request int32, arg any) {
		opus.EncoderCtl(tls, enc, request, opus.VaList(va, arg))
	}

	if !decodeOnly {
		enc = opus.EncoderCreate(tls, samplingRate, channels, application, errPtr)
		if *errCode != 0 {
			fmt.Fprintf(os.Stderr, "Cannot create encoder: %s\n", goString(opus.Strerror(tls, *errCode)))
			return exitFailure
		}
		encCtl(setBitrateRequest, bitrateBps)
		encCtl(setBandwidthRequest, bandwidth)
		encCtl(setVBRRequest, useVBR)
		encCtl(setVBRConstraintRequest, cvbr)
		encCtl(setComplexityRequest, complexity)
		encCtl(setInbandFECRequest, useInbandFEC)
		encCtl(setForceChannelsRequest, forcechannels)
		encCtl(setDTXRequest, useDTX)
		encCtl(setPacketLossPercRequest, packetLossPerc)
		encCtl(getLookaheadRequest, skipPtr)
	}
	if !encodeOnly {
		dec = opus.DecoderCreate(tls, samplingRate, channels, errPtr)
		if *errCode != 0 {
			fmt.Fprintf(os.Stderr, "Cannot create decoder: %s\n", goString(opus.Strerror(tls, *errCode)))
			return exitFailure
		}
	}

	var bandwidthString string
	switch bandwidth {
	case bandwidthNarrowband:
		bandwidthString = "narrowband"
	case bandwidthMediumband:
		bandwidthString = "mediumband"
	case bandwidthWideband:
		bandwidthString = "wideband"
	case bandwidthSuperwideband:
		bandwidthString = "superwideband"
	case bandwidthFullband:
		bandwidthString = "fullband"
	case -1000:
		bandwidthString = "auto"
	default:
		bandwidthString = "unknown"
	}

	if decodeOnly {
		fmt.Fprintf(os.Stderr, "Decoding with %d Hz output (%d channels)\n", int64(samplingRate), channels)
	} else {
		fmt.Fprintf(os.Stderr, "Encoding %d Hz input at %.3f kb/s in %s mode with %d-sample frames.\n",
			int64(samplingRate), float64(bitrateBps)*0.001, bandwidthString, frameSize)
	}

	bufSamples := int(maxFrameSize * channels)
	inPtr := cmalloc(bufSamples * 2)
	defer opus.Free(inPtr)
	outPtr := cmalloc(bufSamples * 2)
	defer opus.Free(outPtr)
	in := unsafe.Slice((*int16)(cptr(inPtr)), bufSamples)
	out := unsafe.Slice((*int16)(cptr(outPtr)), bufSamples)
	fbytes := make([]byte, bufSamples*2)
	var data [2]uintptr
	data[0] = cmalloc(int(maxPayloadBytes))
	defer opus.Free(data[0])
	if useInbandFEC != 0 {
		data[1] = cmalloc(int(maxPayloadBytes))
		defer opus.Free(data[1])
	}
	payload := func(i int32) []byte {
		return unsafe.Slice((*byte)(cptr(data[i])), maxPayloadBytes)
	}
	var ch, intField [4]byte

	for !stop {
		if delayedCelt {
			frameSize = newsize
			delayedCelt = false
		} else if randomFramesize && crand()%20 == 0 {
			newsize = crand() % 6
			switch newsize {
			case 0:
				newsize = samplingRate / 400
			case 1:
				newsize = samplingRate / 200
			case 2:
				newsize = samplingRate / 100
			case 3:
				newsize = samplingRate / 50
			case 4:
				newsize = samplingRate / 25
			case 5:
				newsize = 3 * samplingRate / 50
			}
			for newsize < samplingRate/25 && float64(bitrateBps)-math.Abs(float64(sweepBps)) <= float64(3*12*samplingRate/newsize) {
				newsize = newsize * 2
			}
			if newsize < samplingRate/100 && frameSize >= samplingRate/100 {
				encCtl(setForceModeRequest, int32(modeCeltOnly))
				delayedCelt = true
			} else {
				frameSize = newsize
			}
		}
		if randomFEC && crand()%30 == 0 {
			// OPUS_SET_INBAND_FEC(rand()%4==0) evaluates its argument twice
			// (__opus_check_int), so rand() is called twice as in C.
			_ = crand()
			fec := int32(0)
			if crand()%4 == 0 {
				fec = 1
			}
			encCtl(setInbandFECRequest, fec)
		}
		if decodeOnly {
			if readFull(fin, ch[:]) < 4 {
				break
			}
			lens[toggle] = int32(charToInt(ch[:]))
			if lens[toggle] > maxPayloadBytes || lens[toggle] < 0 {
				fmt.Fprintf(os.Stderr, "Invalid payload length: %d\n", lens[toggle])
				break
			}
			readFull(fin, ch[:])
			encFinalRange[toggle] = charToInt(ch[:])
			n := int32(readFull(fin, payload(toggle)[:lens[toggle]]))
			if n < lens[toggle] {
				fmt.Fprintf(os.Stderr, "Ran out of input, expecting %d bytes got %d\n", lens[toggle], n)
				break
			}
		} else {
			if modeList != nil {
				m := modeList[currMode]
				encCtl(setBandwidthRequest, m.bandwidth)
				encCtl(setForceModeRequest, m.mode)
				encCtl(setForceChannelsRequest, m.channels)
				frameSize = m.frameSize
			}
			frameBytes := int(2 * channels * frameSize)
			currRead = int32(readFull(fin, fbytes[:frameBytes]) / int(2*channels))
			for i := int32(0); i < currRead*channels; i++ {
				in[i] = int16(uint16(fbytes[2*i]) | uint16(fbytes[2*i+1])<<8)
			}
			if currRead < frameSize {
				for i := currRead * channels; i < frameSize*channels; i++ {
					in[i] = 0
				}
				stop = true
			}
			lens[toggle] = opus.Encode(tls, enc, inPtr, frameSize, data[toggle], maxPayloadBytes)
			if sweepBps != 0 {
				bitrateBps += sweepBps
				if sweepMax != 0 {
					if bitrateBps > sweepMax {
						sweepBps = -sweepBps
					} else if bitrateBps < sweepMin {
						sweepBps = -sweepBps
					}
				}
				/* safety */
				if bitrateBps < 1000 {
					bitrateBps = 1000
				}
				encCtl(setBitrateRequest, bitrateBps)
			}
			encCtl(getFinalRangeRequest, encRangePtr+uintptr(toggle)*4)
			if lens[toggle] < 0 {
				fmt.Fprintf(os.Stderr, "opus_encode() returned %d\n", lens[toggle])
				return exitFailure
			}
			currModeCount += frameSize
			if currModeCount > modeSwitchTime && currMode < nbModesInList-1 {
				currMode++
				currModeCount = 0
			}
		}

		if encodeOnly {
			intToChar(uint32(lens[toggle]), intField[:])
			if !fwrite(intField[:]) {
				fmt.Fprint(os.Stderr, "Error writing.\n")
				return exitFailure
			}
			intToChar(encFinalRange[toggle], intField[:])
			if !fwrite(intField[:]) {
				fmt.Fprint(os.Stderr, "Error writing.\n")
				return exitFailure
			}
			if !fwrite(payload(toggle)[:lens[toggle]]) {
				fmt.Fprint(os.Stderr, "Error writing.\n")
				return exitFailure
			}
		} else {
			lost = lens[toggle] == 0 || packetLossPerc > 0 && crand()%100 < packetLossPerc
			if count >= useInbandFEC {
				/* delay by one packet when using in-band FEC */
				if useInbandFEC != 0 {
					if lostPrev {
						/* attempt to decode with in-band FEC from next packet */
						p := data[toggle]
						if lost {
							p = 0
						}
						outputSamples = opus.Decode(tls, dec, p, lens[toggle], outPtr, maxFrameSize, 1)
					} else {
						/* regular decode */
						outputSamples = opus.Decode(tls, dec, data[1-toggle], lens[1-toggle], outPtr, maxFrameSize, 0)
					}
				} else {
					p := data[toggle]
					if lost {
						p = 0
					}
					outputSamples = opus.Decode(tls, dec, p, lens[toggle], outPtr, maxFrameSize, 0)
				}
				if outputSamples > 0 {
					if outputSamples > *skip {
						n := (outputSamples - *skip) * channels
						for i := int32(0); i < n; i++ {
							s := out[i+*skip*channels]
							fbytes[2*i] = byte(s)
							fbytes[2*i+1] = byte(s >> 8)
						}
						if !fwrite(fbytes[:2*n]) {
							fmt.Fprint(os.Stderr, "Error writing.\n")
							return exitFailure
						}
					}
					if outputSamples < *skip {
						*skip -= outputSamples
					} else {
						*skip = 0
					}
				} else {
					fmt.Fprintf(os.Stderr, "error decoding frame: %s\n", goString(opus.Strerror(tls, outputSamples)))
				}
			}
		}

		if !encodeOnly {
			opus.DecoderCtl(tls, dec, getFinalRangeRequest, opus.VaList(va, decRangePtr))
		}
		/* compare final range encoder rng values of encoder and decoder */
		if r := encFinalRange[toggle^useInbandFEC]; r != 0 && !encodeOnly && !lost && !lostPrev && *decFinalRange != r {
			fmt.Fprintf(os.Stderr, "Error: Range coder state mismatch between encoder and decoder in frame %d: 0x%8x vs 0x%8x\n",
				int64(count), uint64(r), uint64(*decFinalRange))
			return exitFailure
		}

		lostPrev = lost

		/* count bits */
		bits += float64(lens[toggle] * 8)
		bitsMax = max(float64(lens[toggle]*8), bitsMax)
		if count >= useInbandFEC {
			nrg = 0.0
			if !decodeOnly {
				for k := int32(0); k < frameSize*channels; k++ {
					nrg = nrg + float64(float64(in[k])*float64(in[k]))
				}
			}
			if nrg/float64(frameSize*channels) > 1e5 {
				bitsAct += float64(lens[toggle] * 8)
				countAct++
			}
			/* Variance */
			bits2 += float64(lens[toggle] * lens[toggle] * 64)
		}
		count++
		toggle = (toggle + useInbandFEC) & 1
	}
	// The explicit float64 conversions keep each product rounded separately
	// (no fused multiply-add), matching the reference arithmetic.
	fmt.Fprintf(os.Stderr, "average bitrate:             %7.3f kb/s\n",
		float64(float64(0.001*bits)*float64(samplingRate))/float64(float64(frameSize)*float64(count)))
	fmt.Fprintf(os.Stderr, "maximum bitrate:             %7.3f bkp/s\n",
		float64(float64(0.001*bitsMax)*float64(samplingRate))/float64(frameSize))
	if !decodeOnly {
		fmt.Fprintf(os.Stderr, "active bitrate:              %7.3f kb/s\n",
			float64(float64(0.001*bitsAct)*float64(samplingRate))/float64(float64(frameSize)*float64(countAct)))
	}
	fmt.Fprintf(os.Stderr, "bitrate standard deviation:  %7.3f kb/s\n",
		float64(float64(0.001*math.Sqrt(bits2/float64(count)-float64(bits*bits)/float64(float64(count)*float64(count))))*float64(samplingRate))/float64(frameSize))
	/* Close any files to which intermediate results were stored */
	opus.EncoderDestroy(tls, enc)
	opus.DecoderDestroy(tls, dec)
	return exitSuccess
}
