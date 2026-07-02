// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.

package opus

import "testing"

func TestReferenceDemoReturnsOnArgumentErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"empty", nil},
		{"usage", []string{"opus_demo"}},
		{"encoder option while decoding", []string{"opus_demo", "-d", "48000", "2", "-cbr", "input", "output"}},
		{"decoder option while encoding", []string{"opus_demo", "-e", "audio", "48000", "2", "64000", "-loss", "10", "input", "output"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if code := RunReferenceDemo(tc.args); code != 1 {
				t.Fatalf("exit code = %d, want 1", code)
			}
		})
	}
}
