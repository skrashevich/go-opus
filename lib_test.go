// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.

package opus

import (
	"math"
	"runtime"
	"testing"
	"unsafe"

	"github.com/skrashevich/go-opus/internal/libc"
)

const (
	testSampleRate = 48000
	testChannels   = 2
	testFrameSize  = testSampleRate / 50
	testMaxPacket  = 4096
)

func newTestTLS(tb testing.TB) *libc.TLS {
	tb.Helper()

	tls := libc.NewTLS()
	tb.Cleanup(tls.Close)
	return tls
}

func makeTestPCM(frameSize, channels int) []opus_int16 {
	pcm := make([]opus_int16, frameSize*channels)
	for i := range frameSize {
		left := int16(math.Sin(2*math.Pi*440*float64(i)/testSampleRate) * 12000)
		right := int16(math.Sin(2*math.Pi*660*float64(i)/testSampleRate) * 10000)
		pcm[i*channels] = left
		if channels > 1 {
			pcm[i*channels+1] = right
		}
	}
	return pcm
}

func ptrOf[T any](v []T) uintptr {
	if len(v) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&v[0]))
}

func testCString(p uintptr) string {
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

func newEncoder(tb testing.TB, tls *libc.TLS) uintptr {
	tb.Helper()

	var err int32
	enc := opus_encoder_create(tls, testSampleRate, testChannels, OPUS_APPLICATION_AUDIO, uintptr(unsafe.Pointer(&err)))
	if enc == 0 || err != OPUS_OK {
		tb.Fatalf("opus_encoder_create() = %#x, err = %d", enc, err)
	}
	tb.Cleanup(func() { opus_encoder_destroy(tls, enc) })
	return enc
}

func newDecoder(tb testing.TB, tls *libc.TLS) uintptr {
	tb.Helper()

	var err int32
	dec := opus_decoder_create(tls, testSampleRate, testChannels, uintptr(unsafe.Pointer(&err)))
	if dec == 0 || err != OPUS_OK {
		tb.Fatalf("opus_decoder_create() = %#x, err = %d", dec, err)
	}
	tb.Cleanup(func() { opus_decoder_destroy(tls, dec) })
	return dec
}

func encodeTestFrame(tb testing.TB, tls *libc.TLS, enc uintptr, pcm []opus_int16) []byte {
	tb.Helper()

	packet := make([]byte, testMaxPacket)
	n := opus_encode(tls, enc, ptrOf(pcm), testFrameSize, ptrOf(packet), int32(len(packet)))
	runtime.KeepAlive(pcm)
	if n <= 0 {
		tb.Fatalf("opus_encode() = %d", n)
	}
	return packet[:n]
}

func TestCodecCreateRejectsInvalidParameters(t *testing.T) {
	tls := newTestTLS(t)

	var err int32
	if enc := opus_encoder_create(tls, 44100, testChannels, OPUS_APPLICATION_AUDIO, uintptr(unsafe.Pointer(&err))); enc != 0 {
		opus_encoder_destroy(tls, enc)
		t.Fatalf("opus_encoder_create() unexpectedly succeeded for 44100 Hz")
	}
	if err != -1 {
		t.Fatalf("encoder error = %d, want -1", err)
	}

	err = 0
	if dec := opus_decoder_create(tls, testSampleRate, 3, uintptr(unsafe.Pointer(&err))); dec != 0 {
		opus_decoder_destroy(tls, dec)
		t.Fatalf("opus_decoder_create() unexpectedly succeeded for 3 channels")
	}
	if err != -1 {
		t.Fatalf("decoder error = %d, want -1", err)
	}
}

func TestCodecSizesAndVersion(t *testing.T) {
	tls := newTestTLS(t)

	if got := opus_encoder_get_size(tls, 1); got <= 0 {
		t.Fatalf("opus_encoder_get_size(1) = %d", got)
	}
	if got := opus_encoder_get_size(tls, 2); got <= 0 {
		t.Fatalf("opus_encoder_get_size(2) = %d", got)
	}
	if got := opus_encoder_get_size(tls, 3); got != 0 {
		t.Fatalf("opus_encoder_get_size(3) = %d, want 0", got)
	}
	if got := opus_decoder_get_size(tls, 1); got <= 0 {
		t.Fatalf("opus_decoder_get_size(1) = %d", got)
	}
	if got := opus_decoder_get_size(tls, 2); got <= 0 {
		t.Fatalf("opus_decoder_get_size(2) = %d", got)
	}
	if got := opus_decoder_get_size(tls, 0); got != 0 {
		t.Fatalf("opus_decoder_get_size(0) = %d, want 0", got)
	}

	if got := testCString(opus_get_version_string(tls)); got != "libopus 1.0.0" {
		t.Fatalf("opus_get_version_string() = %q", got)
	}
	if got := testCString(opus_strerror(tls, OPUS_OK)); got != "success" {
		t.Fatalf("opus_strerror(OPUS_OK) = %q", got)
	}
}

func TestEncodeDecodeRoundTripPacket(t *testing.T) {
	tls := newTestTLS(t)
	enc := newEncoder(t, tls)
	dec := newDecoder(t, tls)
	pcm := makeTestPCM(testFrameSize, testChannels)

	packet := encodeTestFrame(t, tls, enc, pcm)
	packetPtr := ptrOf(packet)
	if got := opus_packet_get_nb_frames(tls, packetPtr, int32(len(packet))); got != 1 {
		t.Fatalf("opus_packet_get_nb_frames() = %d, want 1", got)
	}
	if got := opus_packet_get_nb_channels(tls, packetPtr); got != testChannels {
		t.Fatalf("opus_packet_get_nb_channels() = %d, want %d", got, testChannels)
	}
	if got := opus_packet_get_samples_per_frame(tls, packetPtr, testSampleRate); got != testFrameSize {
		t.Fatalf("opus_packet_get_samples_per_frame() = %d, want %d", got, testFrameSize)
	}
	if got := opus_decoder_get_nb_samples(tls, dec, packetPtr, int32(len(packet))); got != testFrameSize {
		t.Fatalf("opus_decoder_get_nb_samples() = %d, want %d", got, testFrameSize)
	}

	decoded := make([]opus_int16, testFrameSize*testChannels)
	n := opus_decode(tls, dec, packetPtr, int32(len(packet)), ptrOf(decoded), testFrameSize, 0)
	runtime.KeepAlive(packet)
	if n != testFrameSize {
		t.Fatalf("opus_decode() = %d, want %d", n, testFrameSize)
	}
}

func BenchmarkEncode20msStereo48k(b *testing.B) {
	tls := newTestTLS(b)
	enc := newEncoder(b, tls)
	pcm := makeTestPCM(testFrameSize, testChannels)
	packet := make([]byte, testMaxPacket)
	pcmPtr := ptrOf(pcm)
	packetPtr := ptrOf(packet)

	b.SetBytes(int64(len(pcm) * 2))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if n := opus_encode(tls, enc, pcmPtr, testFrameSize, packetPtr, int32(len(packet))); n <= 0 {
			b.Fatalf("opus_encode() = %d", n)
		}
	}
	runtime.KeepAlive(pcm)
	runtime.KeepAlive(packet)
}

func BenchmarkDecode20msStereo48k(b *testing.B) {
	tls := newTestTLS(b)
	enc := newEncoder(b, tls)
	dec := newDecoder(b, tls)
	packet := encodeTestFrame(b, tls, enc, makeTestPCM(testFrameSize, testChannels))
	decoded := make([]opus_int16, testFrameSize*testChannels)
	packetPtr := ptrOf(packet)
	decodedPtr := ptrOf(decoded)

	b.SetBytes(int64(len(decoded) * 2))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if n := opus_decode(tls, dec, packetPtr, int32(len(packet)), decodedPtr, testFrameSize, 0); n != testFrameSize {
			b.Fatalf("opus_decode() = %d, want %d", n, testFrameSize)
		}
	}
	runtime.KeepAlive(packet)
	runtime.KeepAlive(decoded)
}

func BenchmarkEncodeDecode20msStereo48k(b *testing.B) {
	tls := newTestTLS(b)
	enc := newEncoder(b, tls)
	dec := newDecoder(b, tls)
	pcm := makeTestPCM(testFrameSize, testChannels)
	packet := make([]byte, testMaxPacket)
	decoded := make([]opus_int16, testFrameSize*testChannels)
	pcmPtr := ptrOf(pcm)
	packetPtr := ptrOf(packet)
	decodedPtr := ptrOf(decoded)

	b.SetBytes(int64(len(pcm) * 2))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := opus_encode(tls, enc, pcmPtr, testFrameSize, packetPtr, int32(len(packet)))
		if n <= 0 {
			b.Fatalf("opus_encode() = %d", n)
		}
		if got := opus_decode(tls, dec, packetPtr, n, decodedPtr, testFrameSize, 0); got != testFrameSize {
			b.Fatalf("opus_decode() = %d, want %d", got, testFrameSize)
		}
	}
	runtime.KeepAlive(pcm)
	runtime.KeepAlive(packet)
	runtime.KeepAlive(decoded)
}
