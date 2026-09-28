// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0
// See LICENSE and NOTICE in the project root.
// https://github.com/skrashevich/go-opus

//go:build !ppc64 && !ppc64le

package opus

// sx16 is int16, so sx16(x) is a plain conversion. On ppc64 it is a function
// that works around a compiler bug; see sx16_ppc64x.go.
type sx16 = int16
