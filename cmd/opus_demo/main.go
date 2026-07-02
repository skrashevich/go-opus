// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.

package main

import (
	"os"
	"runtime/pprof"

	"github.com/skrashevich/go-opus"
)

func main() {
	var profileFile *os.File
	if p := os.Getenv("CPUPROFILE"); p != "" {
		var err error
		profileFile, err = os.Create(p)
		if err == nil {
			_ = pprof.StartCPUProfile(profileFile)
		}
	}

	code := opus.RunReferenceDemo(os.Args)

	if profileFile != nil {
		pprof.StopCPUProfile()
		_ = profileFile.Close()
	}
	os.Exit(code)
}
