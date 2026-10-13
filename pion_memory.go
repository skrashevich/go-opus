// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0

package opus

import (
	"fmt"
	"runtime"
	"unsafe"
)

const (
	codecAllocFail      int32 = -7
	codecBufferTooSmall int32 = -2
	codecInvalidPacket  int32 = -4
)

// The cleanup argument owns C allocations but never refers back to nativeCodec.
// Keeping it separate permits resizing buffers without retaining the cleanup target.
type codecResources struct {
	tls                           *TLS
	state, input, output, scratch uintptr
	inputSize, outputSize         int
}

type nativeCodec struct {
	r       *codecResources
	cleanup runtime.Cleanup
}

func (r *codecResources) close() {
	Free(r.state)
	Free(r.input)
	Free(r.output)
	Free(r.scratch)
	r.tls.Close()
}

func newNativeCodec(rate, channels int, encoder bool) (*nativeCodec, error) {
	r := &codecResources{tls: NewTLS(), scratch: Malloc(128)}
	if r.scratch == 0 {
		r.close()
		return nil, codecError(codecAllocFail)
	}
	if encoder {
		r.state = opus_encoder_create(r.tls, int32(rate), int32(channels), OPUS_APPLICATION_AUDIO, r.scratch)
	} else {
		r.state = opus_decoder_create(r.tls, int32(rate), int32(channels), r.scratch)
	}
	if r.state == 0 {
		code := *(*int32)(unsafe.Pointer(r.scratch))
		r.close()
		return nil, codecError(code)
	}
	n := &nativeCodec{r: r}
	n.cleanup = runtime.AddCleanup(n, (*codecResources).close, r)
	return n, nil
}

func growCodecBuffer(p *uintptr, capacity *int, size int) error {
	if size <= *capacity {
		return nil
	}
	next := Malloc(size)
	if next == 0 {
		return codecError(codecAllocFail)
	}
	Free(*p)
	*p, *capacity = next, size
	return nil
}

func (r *codecResources) buffers(inputBytes, outputBytes int) error {
	if err := growCodecBuffer(&r.input, &r.inputSize, inputBytes); err != nil {
		return err
	}
	return growCodecBuffer(&r.output, &r.outputSize, outputBytes)
}

func (n *nativeCodec) ctl(request, value int32) error {
	defer runtime.KeepAlive(n)
	r := n.r
	return codecError(opus_encoder_ctl(r.tls, r.state, request, VaList(r.scratch+112, value)))
}

func codecError(code int32) error {
	if code >= 0 {
		return nil
	}
	switch code {
	case codecBufferTooSmall:
		return errOutBufferTooSmall
	case codecInvalidPacket:
		return errMalformedPacket
	default:
		return fmt.Errorf("opus: codec error %d", code)
	}
}
