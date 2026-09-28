/* Copyright (c) 2011-2012 Xiph.Org Foundation, Mozilla Corporation
   Written by Jean-Marc Valin and Timothy B. Terriberry */
/*
   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:

   - Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.

   - Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.

   - Neither the name of Internet Society, IETF or IETF Trust, nor the
   names of specific contributors, may be used to endorse or promote
   products derived from this software without specific prior written
   permission.

   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
   ``AS IS'' AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
   LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
   A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER
   OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
   EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
   PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
   PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
   LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
   NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
   SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Command opus_compare is a Go port of opus_compare.c from RFC 6716. It
// compares two 16-bit little-endian PCM files and prints the Opus quality
// metric.
//
// The float32/float64 split and the explicit float32()/float64() conversions
// around products mirror the C code exactly: they keep every intermediate
// rounded as in C and stop the Go compiler from fusing x*y+z into FMA, so the
// printed metric is bit-identical to the reference on every architecture.
package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

const (
	NBANDS        = 21
	NFREQS        = 240
	TEST_WIN_SIZE = 480
	TEST_WIN_STEP = 120
)

func fail(format string, args ...any) int {
	fmt.Fprintf(os.Stderr, format, args...)
	return 1
}

// readPCM16 reads interleaved 16-bit little-endian samples and returns them
// as floats together with the number of sample frames read.
func readPCM16(fin *os.File, nchannels int) ([]float32, int) {
	var buf [1024]byte
	size := 2 * nchannels
	var samples []float32
	nsamples := 0
	for {
		// Like the reference fread (a single read(2) returning n/size),
		// a trailing partial sample frame is dropped.
		n, err := fin.Read(buf[:1024/size*size])
		if err != nil && err != io.EOF {
			// The reference fread returns (size_t)-1/size on a read error,
			// which overflows the subsequent realloc size and panics there.
			panic("invalid malloc size")
		}
		nread := n / size
		if nread <= 0 {
			break
		}
		for xi := 0; xi < nread; xi++ {
			for ci := 0; ci < nchannels; ci++ {
				i := 2 * (xi*nchannels + ci)
				s := int16(uint16(buf[i]) | uint16(buf[i+1])<<8)
				samples = append(samples, float32(s))
			}
		}
		nsamples += nread
	}
	if nsamples == 0 {
		// The reference realloc(NULL, 0) returns NULL, reported as OOM.
		fmt.Fprint(os.Stderr, "Out of memory.\n")
		os.Exit(1)
	}
	return samples, nsamples
}

func bandEnergy(out, ps []float32, bands []int, nbands int, in []float32,
	nchannels, nframes, windowSz, step, downsample int) {
	window := make([]float32, windowSz)
	c := make([]float32, windowSz)
	s := make([]float32, windowSz)
	x := make([]float32, nchannels*windowSz)
	psSz := windowSz / 2
	twoPi := float32(2) * pi
	for xj := 0; xj < windowSz; xj++ {
		window[xj] = 0.5 - float32(0.5*float32(math.Cos(float64(twoPi/float32(windowSz-1)*float32(xj)))))
	}
	for xj := 0; xj < windowSz; xj++ {
		c[xj] = float32(math.Cos(float64(twoPi / float32(windowSz) * float32(xj))))
	}
	for xj := 0; xj < windowSz; xj++ {
		s[xj] = float32(math.Sin(float64(twoPi / float32(windowSz) * float32(xj))))
	}
	for xi := 0; xi < nframes; xi++ {
		for ci := 0; ci < nchannels; ci++ {
			for xk := 0; xk < windowSz; xk++ {
				x[ci*windowSz+xk] = window[xk] * in[(xi*step+xk)*nchannels+ci]
			}
		}
		xj := 0
		for bi := 0; bi < nbands; bi++ {
			var p [2]float32
			for ; xj < bands[bi+1]; xj++ {
				for ci := 0; ci < nchannels; ci++ {
					var re, im float32
					ti := 0
					for xk := 0; xk < windowSz; xk++ {
						re = re + float32(c[ti]*x[ci*windowSz+xk])
						im = im - float32(s[ti]*x[ci*windowSz+xk])
						ti += xj
						if ti >= windowSz {
							ti -= windowSz
						}
					}
					re = re * float32(downsample)
					im = im * float32(downsample)
					k := (xi*psSz+xj)*nchannels + ci
					ps[k] = float32(re*re) + float32(im*im) + 100000
					p[ci] += ps[k]
				}
			}
			if out != nil {
				out[(xi*nbands+bi)*nchannels] = p[0] / float32(bands[bi+1]-bands[bi])
				if nchannels == 2 {
					out[(xi*nbands+bi)*nchannels+1] = p[1] / float32(bands[bi+1]-bands[bi])
				}
			}
		}
	}
}

// pi is a variable so 2*pi is computed in float32 exactly as in C.
var pi float32 = 3.14159265

/*Bands on which we compute the pseudo-NMR (Bark-derived CELT bands).*/
var BANDS = []int{
	0, 2, 4, 6, 8, 10, 12, 14, 16, 20, 24, 28, 32, 40, 48, 56, 68, 80, 96, 120, 156, 200,
}

// cstr models passing a possibly-NULL C string to fopen/printf("%s"),
// which the reference libc treats as "".
func cstr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// atoi mirrors the reference libc atoi: leading whitespace, optional sign,
// decimal digits accumulated in uint64 (stopping on overflow), truncated to
// int32.
func atoi(str string) int32 {
	i := 0
	neg := false
loop:
	for ; i < len(str); i++ {
		switch str[i] {
		case ' ', '\t', '\n', '\r', '\v', '\f':
		case '+':
			i++
			break loop
		case '-':
			i++
			neg = true
			break loop
		default:
			break loop
		}
	}
	var n uint64
	for ; i < len(str) && str[i] >= '0' && str[i] <= '9'; i++ {
		n0 := n
		n = 10*n + uint64(str[i]-'0')
		if n < n0 {
			n = n0
			break
		}
	}
	if neg {
		return int32(-n)
	}
	return int32(n)
}

func run() int {
	argc := len(os.Args)
	// argv mirrors C's NULL-terminated argv array.
	argv := make([]*string, argc+1)
	for i := range os.Args {
		argv[i] = &os.Args[i]
	}
	if argc < 3 || argc > 6 {
		return fail("Usage: %s [-s] [-r rate2] <file1.sw> <file2.sw>\n", *argv[0])
	}
	nchannels := 1
	if *argv[1] == "-s" {
		nchannels = 2
		argv = argv[1:]
	}
	rate := uint32(48000)
	ybands := NBANDS
	yfreqs := NFREQS
	downsample := 1
	if *argv[1] == "-r" {
		rate = uint32(atoi(*argv[2]))
		if rate != 8000 && rate != 12000 && rate != 16000 && rate != 24000 && rate != 48000 {
			return fail("Sampling rate must be 8000, 12000, 16000, 24000, or 48000\n")
		}
		downsample = int(48000 / rate)
		switch rate {
		case 8000:
			ybands = 13
		case 12000:
			ybands = 15
		case 16000:
			ybands = 17
		case 24000:
			ybands = 19
		}
		yfreqs = NFREQS / downsample
		argv = argv[2:]
	}
	fin1, err := os.Open(cstr(argv[1]))
	if err != nil {
		return fail("Error opening '%s'.\n", cstr(argv[1]))
	}
	fin2, err := os.Open(cstr(argv[2]))
	if err != nil {
		fin1.Close()
		return fail("Error opening '%s'.\n", cstr(argv[2]))
	}
	/*Read in the data and allocate scratch space.*/
	x, xlength := readPCM16(fin1, 2)
	if nchannels == 1 {
		for xi := 0; xi < xlength; xi++ {
			x[xi] = float32(0.5 * float64(x[2*xi]+x[2*xi+1]))
		}
	}
	fin1.Close()
	y, ylength := readPCM16(fin2, nchannels)
	fin2.Close()
	if xlength != ylength*downsample {
		return fail("Sample counts do not match (%d!=%d).\n", xlength, ylength*downsample)
	}
	if xlength < TEST_WIN_SIZE {
		return fail("Insufficient sample data (%d<%d).\n", xlength, TEST_WIN_SIZE)
	}
	nframes := (xlength - TEST_WIN_SIZE + TEST_WIN_STEP) / TEST_WIN_STEP
	xb := make([]float32, nframes*NBANDS*nchannels)
	X := make([]float32, nframes*NFREQS*nchannels)
	Y := make([]float32, nframes*yfreqs*nchannels)
	/*Compute the per-band spectral energy of the original signal
	  and the error.*/
	bandEnergy(xb, X, BANDS, NBANDS, x, nchannels, nframes,
		TEST_WIN_SIZE, TEST_WIN_STEP, 1)
	bandEnergy(nil, Y, BANDS, ybands, y, nchannels, nframes,
		TEST_WIN_SIZE/downsample, TEST_WIN_STEP/downsample, downsample)
	for xi := 0; xi < nframes; xi++ {
		/*Frequency masking (low to high): 10 dB/Bark slope.*/
		for bi := 1; bi < NBANDS; bi++ {
			for ci := 0; ci < nchannels; ci++ {
				xb[(xi*NBANDS+bi)*nchannels+ci] += float32(0.1 * xb[(xi*NBANDS+bi-1)*nchannels+ci])
			}
		}
		/*Frequency masking (high to low): 15 dB/Bark slope.*/
		for bi := NBANDS - 2; bi >= 0; bi-- {
			for ci := 0; ci < nchannels; ci++ {
				xb[(xi*NBANDS+bi)*nchannels+ci] += float32(0.03 * xb[(xi*NBANDS+bi+1)*nchannels+ci])
			}
		}
		if xi > 0 {
			/*Temporal masking: -3 dB/2.5ms slope.*/
			for bi := 0; bi < NBANDS; bi++ {
				for ci := 0; ci < nchannels; ci++ {
					xb[(xi*NBANDS+bi)*nchannels+ci] += float32(0.5 * xb[((xi-1)*NBANDS+bi)*nchannels+ci])
				}
			}
		}
		/* Allowing some cross-talk */
		if nchannels == 2 {
			for bi := 0; bi < NBANDS; bi++ {
				l := xb[(xi*NBANDS+bi)*nchannels+0]
				r := xb[(xi*NBANDS+bi)*nchannels+1]
				xb[(xi*NBANDS+bi)*nchannels+0] += float32(0.01 * r)
				xb[(xi*NBANDS+bi)*nchannels+1] += float32(0.01 * l)
			}
		}
		/* Apply masking */
		for bi := 0; bi < ybands; bi++ {
			for xj := BANDS[bi]; xj < BANDS[bi+1]; xj++ {
				for ci := 0; ci < nchannels; ci++ {
					X[(xi*NFREQS+xj)*nchannels+ci] += float32(0.1 * xb[(xi*NBANDS+bi)*nchannels+ci])
					Y[(xi*yfreqs+xj)*nchannels+ci] += float32(0.1 * xb[(xi*NBANDS+bi)*nchannels+ci])
				}
			}
		}
	}

	/* Average of consecutive frames to make comparison slightly less sensitive */
	for bi := 0; bi < ybands; bi++ {
		for xj := BANDS[bi]; xj < BANDS[bi+1]; xj++ {
			for ci := 0; ci < nchannels; ci++ {
				xtmp := X[xj*nchannels+ci]
				ytmp := Y[xj*nchannels+ci]
				for xi := 1; xi < nframes; xi++ {
					xtmp2 := X[(xi*NFREQS+xj)*nchannels+ci]
					ytmp2 := Y[(xi*yfreqs+xj)*nchannels+ci]
					X[(xi*NFREQS+xj)*nchannels+ci] += xtmp
					Y[(xi*yfreqs+xj)*nchannels+ci] += ytmp
					xtmp = xtmp2
					ytmp = ytmp2
				}
			}
		}
	}

	/*If working at a lower sampling rate, don't take into account the last
	   300 Hz to allow for different transition bands.
	  For 12 kHz, we don't skip anything, because the last band already skips
	   400 Hz.*/
	var maxCompare int
	if rate == 48000 {
		maxCompare = BANDS[NBANDS]
	} else if rate == 12000 {
		maxCompare = BANDS[ybands]
	} else {
		maxCompare = BANDS[ybands] - 3
	}
	errW := 0.0
	for xi := 0; xi < nframes; xi++ {
		Ef := 0.0
		for bi := 0; bi < ybands; bi++ {
			Eb := 0.0
			for xj := BANDS[bi]; xj < BANDS[bi+1] && xj < maxCompare; xj++ {
				for ci := 0; ci < nchannels; ci++ {
					re := Y[(xi*yfreqs+xj)*nchannels+ci] / X[(xi*NFREQS+xj)*nchannels+ci]
					im := float32(float64(re) - math.Log(float64(re)) - 1)
					/*Make comparison less sensitive around the SILK/CELT cross-over to
					  allow for mode freedom in the filters.*/
					if xj >= 79 && xj <= 81 {
						im *= 0.1
					}
					if xj == 80 {
						im *= 0.1
					}
					Eb += float64(im)
				}
			}
			Eb /= float64((BANDS[bi+1] - BANDS[bi]) * nchannels)
			Ef += float64(Eb * Eb)
		}
		/*Using a fixed normalization value means we're willing to accept slightly
		  lower quality for lower sampling rates.*/
		Ef /= NBANDS
		Ef *= Ef
		errW += float64(Ef * Ef)
	}
	errW = math.Pow(errW/float64(nframes), 1.0/16)
	Q := float32(100 * (1 - float64(0.5*math.Log(1+errW))/math.Log(1.13)))
	if Q < 0 {
		fmt.Fprint(os.Stderr, "Test vector FAILS\n")
		return fail("Internal weighted error is %s\n", cfmt("%f", errW))
	}
	fmt.Fprint(os.Stderr, "Test vector PASSES\n")
	fmt.Fprintf(os.Stderr, "Opus quality metric: %s %% (internal weighted error is %s)\n",
		cfmt("%.1f", float64(Q)), cfmt("%f", errW))
	return 0
}

// cfmt formats a double like C's %f, spelling NaN and infinities as C does.
func cfmt(verb string, v float64) string {
	switch {
	case math.IsNaN(v):
		return "nan"
	case math.IsInf(v, 0):
		return "inf"
	}
	return fmt.Sprintf(verb, v)
}

func main() {
	os.Exit(run())
}
