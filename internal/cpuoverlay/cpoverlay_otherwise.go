//go:build !arm64 || (!darwin && !android)

package cpuoverlay

import "golang.org/x/sys/cpu"

// This file is built when we're not on arm64. In which case we can
// just return the (false) values hold by cpu.ARM64.HasXXX.

// arm64HasAES returns whether the CPU supports AES.
func arm64HasAES() bool {
	return cpu.ARM64.HasAES
}

// arm64HasPMULL returns whether the CPU supports PMULL.
func arm64HasPMULL() bool {
	return cpu.ARM64.HasPMULL
}
