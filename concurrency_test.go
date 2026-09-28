// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

package opus

import (
	"bytes"
	"sync"
	"testing"
	"unsafe"

	"github.com/skrashevich/go-opus/internal/libc"
)

func encodeStream(t *testing.T) []byte {
	tls := libc.NewTLS()
	defer tls.Close()
	const frame, ch = 960, 2
	status := libc.Xmalloc(tls, 4)
	pcm := libc.Xmalloc(tls, frame*ch*2)
	pkt := libc.Xmalloc(tls, 4000)
	defer libc.Xfree(tls, status)
	defer libc.Xfree(tls, pcm)
	defer libc.Xfree(tls, pkt)
	enc := opus_encoder_create(tls, 48000, ch, OPUS_APPLICATION_AUDIO, status)
	defer opus_encoder_destroy(tls, enc)
	s := unsafe.Slice((*int16)(unsafe.Pointer(pcm)), frame*ch)
	var out []byte
	for i := 0; i < 200; i++ {
		for j := range s {
			s[j] = int16((i*977 + j*131) % 20000)
		}
		n := opus_encode(tls, enc, pcm, frame, pkt, 4000)
		out = append(out, unsafe.Slice((*byte)(unsafe.Pointer(pkt)), n)...)
	}
	return out
}

// TestConcurrentEncoders runs independent encoders on separate goroutines;
// each must produce the same stream as a serial run.
func TestConcurrentEncoders(t *testing.T) {
	want := encodeStream(t)
	var wg sync.WaitGroup
	bad := make([]bool, 8)
	for g := range bad {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bad[g] = !bytes.Equal(encodeStream(t), want)
		}()
	}
	wg.Wait()
	for g, b := range bad {
		if b {
			t.Errorf("goroutine %d: output differs from serial run", g)
		}
	}
}
