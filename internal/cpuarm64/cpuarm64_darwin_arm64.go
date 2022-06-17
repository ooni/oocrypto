// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm64 && darwin && !ios

package cpuarm64

// There are no hw.optional sysctl values for the below features on Mac OS 11.0
// to detect their supported state dynamically. Assume the CPU features that
// Apple Silicon M1 supports to be available as a minimal set of features
// to all Go programs running on darwin/arm64.

// HasAES returns whether the CPU supports AES.
func HasAES() bool {
	return true
}

// HasPMULL returns whether the CPU supports PMULL.
func HasPMULL() bool {
	return true
}
