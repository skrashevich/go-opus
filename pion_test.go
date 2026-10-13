// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0

package opus

import (
	"encoding/binary"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// Compile independent consumers under both module identities, without access
// to a second copy of go-opus or to the module network.
func TestPionModuleReplacement(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	source, err := os.ReadFile(filepath.Join(root, "testdata", "pion-client", "client_test.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, module := range []string{"github.com/pion/opus", "github.com/skrashevich/go-opus"} {
		t.Run(module, func(t *testing.T) {
			client := t.TempDir()
			mod := "module pion-client\n\ngo 1.26.1\n\nrequire " + module + " v0.1.0\n\n" +
				"replace " + module + " => " + strconv.Quote(root) + "\n"
			if err := os.WriteFile(filepath.Join(client, "go.mod"), []byte(mod), 0600); err != nil {
				t.Fatal(err)
			}
			clientSource := strings.ReplaceAll(string(source), "github.com/pion/opus", module)
			if err := os.WriteFile(filepath.Join(client, "client_test.go"), []byte(clientSource), 0600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.CommandContext(t.Context(), "go", "test", "-mod=mod", "./...")
			cmd.Dir = client
			cmd.Env = append(os.Environ(), "GOWORK=off", "GOPROXY=off", "GOSUMDB=off")
			tidy := exec.CommandContext(t.Context(), "go", "mod", "tidy")
			tidy.Dir, tidy.Env = cmd.Dir, cmd.Env
			if out, err := tidy.CombinedOutput(); err != nil {
				t.Fatalf("consumer tidy: %v\n%s", err, out)
			}
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("consumer: %v\n%s", err, out)
			}
		})
	}
}

func TestPionRoundTrip(t *testing.T) {
	for _, channels := range []int{1, 2} {
		e, err := NewEncoder(WithChannels(channels), WithBandwidth(BandwidthFullband))
		if err != nil {
			t.Fatal(err)
		}
		pcm := make([]float32, 960*channels)
		for i := range pcm {
			pcm[i] = float32(.25 * math.Sin(float64(i/channels)*.1))
		}
		packet := make([]byte, 1500)
		n, err := e.EncodeFloat32(pcm, packet)
		if err != nil {
			t.Fatal(err)
		}
		packet = packet[:n]
		for _, rate := range []int{8000, 12000, 16000, 24000, 48000} {
			for _, outputChannels := range []int{1, 2} {
				d, err := NewDecoderWithOutput(rate, outputChannels)
				if err != nil {
					t.Fatal(err)
				}
				out := make([]int16, rate/50*outputChannels+3)
				out[len(out)-1] = 1234
				count, err := d.DecodeToInt16(packet, out)
				if err != nil || count != rate/50 {
					t.Fatalf("%d/%d: %d, %v", rate, outputChannels, count, err)
				}
				if out[len(out)-1] != 1234 {
					t.Fatal("output tail overwritten")
				}
				if err := d.Init(rate, outputChannels); err != nil {
					t.Fatal(err)
				}
				bytes := make([]byte, count*outputChannels*2)
				bw, stereo, err := d.Decode(packet, bytes)
				if err != nil || bw != BandwidthFullband || stereo != (channels == 2) {
					t.Fatalf("metadata: %v %v %v", bw, stereo, err)
				}
				for i := range count * outputChannels {
					if int16(binary.LittleEndian.Uint16(bytes[2*i:])) != out[i] {
						t.Fatal("S16LE differs from int16")
					}
				}
			}
		}
	}
}

func TestPionInvalidInputPreservesState(t *testing.T) {
	e, err := NewEncoder()
	if err != nil {
		t.Fatal(err)
	}
	packet := make([]byte, 1500)
	n, err := e.Encode(make([]byte, 1920), packet)
	if err != nil {
		t.Fatal(err)
	}
	packet = packet[:n]
	d := NewDecoder()
	if _, err := d.DecodeToFloat32(nil, make([]float32, 960)); err == nil {
		t.Fatal("empty packet accepted")
	}
	if _, err := d.DecodeToInt16(packet, make([]int16, 959)); err == nil {
		t.Fatal("short output accepted")
	}
	if err := d.Init(44100, 1); err == nil {
		t.Fatal("invalid rate accepted")
	}
	got, want := make([]float32, 960), make([]float32, 960)
	fresh := NewDecoder()
	if _, err := d.DecodeToFloat32(packet, got); err != nil {
		t.Fatal(err)
	}
	if _, err := fresh.DecodeToFloat32(packet, want); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, want) {
		t.Fatal("failed call advanced decoder state")
	}
	if _, err := e.Encode([]byte{0}, packet); err == nil {
		t.Fatal("odd input accepted")
	}
	if _, err := e.EncodeFloat32(make([]float32, 959), packet); err == nil {
		t.Fatal("invalid frame accepted")
	}
	if _, err := e.EncodeFloat32(make([]float32, 960), nil); err == nil {
		t.Fatal("short encoder output accepted")
	}
}

func TestPionSILKAndPLC(t *testing.T) {
	e, err := NewEncoder()
	if err != nil {
		t.Fatal(err)
	}
	for _, bw := range []Bandwidth{BandwidthNarrowband, BandwidthMediumband, BandwidthWideband} {
		for _, units := range []int{1, 2, 3} {
			packet := make([]byte, 1500)
			n, err := e.EncodeSILK(make([]int16, bw.SampleRate()/50*units), bw, packet)
			if err != nil {
				t.Fatal(err)
			}
			d := NewDecoder()
			out := make([]int16, 960*units)
			count, err := d.DecodeToInt16(packet[:n], out)
			if err != nil || count != len(out) {
				t.Fatalf("SILK %v/%d: %d, %v", bw, units, count, err)
			}
			if err := d.DecodePLC(out); err != nil {
				t.Fatal(err)
			}
		}
	}
	d := NewDecoder()
	out := make([]int16, 960)
	for i := range out {
		out[i] = 123
	}
	if err := d.DecodePLC(out); err != nil {
		t.Fatal(err)
	}
	for _, v := range out {
		if v != 0 {
			t.Fatal("unprimed PLC is not silence")
		}
	}
	if err := d.DecodePLC(make([]int16, 480)); err == nil {
		t.Fatal("unprimed short PLC accepted")
	}
}

func TestPionEncoderOptions(t *testing.T) {
	e, err := NewEncoder(WithVBR(true), WithConstrainedVBR(false), WithApplication(ApplicationVoIP), WithComplexity(7), WithMaxBandwidth(BandwidthWideband))
	if err != nil {
		t.Fatal(err)
	}
	if !e.VBR() || e.ConstrainedVBR() || e.Application() != ApplicationVoIP || e.Complexity() != 7 || e.MaxBandwidth() != BandwidthWideband || e.Bandwidth() != BandwidthAuto {
		t.Fatal("options lost")
	}
	for _, opt := range []EncoderOption{WithSampleRate(44100), WithChannels(3), WithBitrate(5999), WithComplexity(11), WithApplication(0), WithBandwidth(BandwidthAuto), WithMaxBandwidth(BandwidthMediumband)} {
		if _, err := NewEncoder(opt); err == nil {
			t.Fatal("invalid option accepted")
		}
	}
	if err := e.SetLossRate(101); err == nil {
		t.Fatal("invalid loss rate accepted")
	}
	if err := e.SetLossRate(12); err != nil || e.LossRate() != 12 {
		t.Fatal("loss rate lost")
	}
	if err := e.SetBandwidth(BandwidthNarrowband); err != nil {
		t.Fatal(err)
	}
	packet := make([]byte, 1500)
	n, err := e.EncodeFloat32(make([]float32, 960), packet)
	if err != nil {
		t.Fatal(err)
	}
	d := NewDecoder()
	bw, _, err := d.DecodeFloat32(packet[:n], make([]float32, 960))
	if err != nil || bw != BandwidthNarrowband {
		t.Fatalf("bandwidth ctl: %v %v", bw, err)
	}
}
