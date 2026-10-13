// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0

package opus_test

import (
	"github.com/pion/opus"
	"github.com/pion/opus/pkg/oggreader"
	"testing"
)

// These assignments protect the complete public codec method signatures from
// upstream 44637de087b3f506170ab738ed96ac1be338006b (plus v0.1.0 Configuration).
var (
	_ func() opus.Decoder                                = opus.NewDecoder
	_ func(int, int) (opus.Decoder, error)               = opus.NewDecoderWithOutput
	_ func(...opus.EncoderOption) (*opus.Encoder, error) = opus.NewEncoder
	_ interface {
		Init(int, int) error
		Decode([]byte, []byte) (opus.Bandwidth, bool, error)
		DecodeFloat32([]byte, []float32) (opus.Bandwidth, bool, error)
		DecodeToInt16([]byte, []int16) (int, error)
		DecodeToFloat32([]byte, []float32) (int, error)
		DecodePLC([]int16) error
	} = (*opus.Decoder)(nil)
	_ interface {
		Encode([]byte, []byte) (int, error)
		EncodeFloat32([]float32, []byte) (int, error)
		EncodeSILK([]int16, opus.Bandwidth, []byte) (int, error)
		SetBitrate(int) error
		SetComplexity(int) error
		SetApplication(opus.Application) error
		SetVBR(bool)
		SetConstrainedVBR(bool)
		SetLossRate(int) error
		SetBandwidth(opus.Bandwidth) error
		SetMaxBandwidth(opus.Bandwidth) error
		Application() opus.Application
		VBR() bool
		Complexity() int
		ConstrainedVBR() bool
		LossRate() int
		Bandwidth() opus.Bandwidth
		MaxBandwidth() opus.Bandwidth
	} = (*opus.Encoder)(nil)
	_ opus.Configuration = 31
	_ *oggreader.OggReader
)

func TestOriginalPionImports(t *testing.T) {
	e, err := opus.NewEncoder(opus.WithSampleRate(48000), opus.WithChannels(2), opus.WithBitrate(64000), opus.WithComplexity(5), opus.WithApplication(opus.ApplicationAudio), opus.WithVBR(true), opus.WithConstrainedVBR(true), opus.WithBandwidth(opus.BandwidthFullband), opus.WithMaxBandwidth(opus.BandwidthFullband))
	if err != nil {
		t.Fatal(err)
	}
	packet := make([]byte, 1500)
	n, err := e.EncodeFloat32(make([]float32, 1920), packet)
	if err != nil {
		t.Fatal(err)
	}
	d, err := opus.NewDecoderWithOutput(48000, 2)
	if err != nil {
		t.Fatal(err)
	}
	count, err := d.DecodeToInt16(packet[:n], make([]int16, 1920))
	if err != nil || count != 960 {
		t.Fatalf("decode: %d %v", count, err)
	}
}
