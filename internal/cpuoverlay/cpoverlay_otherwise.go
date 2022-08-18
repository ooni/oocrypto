//go:build !arm64 || (!darwin && !android)

package cpuoverlay

import "golang.org/x/sys/cpu"

// arm64HasAES returns whether the CPU supports AES.
func arm64HasAES() bool {
	return cpu.ARM64.HasAES
}

// arm64HasPMULL returns whether the CPU supports PMULL.
func arm64HasPMULL() bool {
	return cpu.ARM64.HasPMULL
}
