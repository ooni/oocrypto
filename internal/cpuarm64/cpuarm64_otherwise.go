//go:build !android || !arm64

package cpuarm64

import "golang.org/x/sys/cpu"

// HasAES returns whether the CPU supports AES.
func HasAES() bool {
	return cpu.ARM64.HasAES
}

// HasPMULL returns whether the CPU supports PMULL.
func HasPMULL() bool {
	return cpu.ARM64.HasPMULL
}
