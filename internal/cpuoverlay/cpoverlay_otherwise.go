//go:build !arm64 || (!darwin && !android)

package cpuoverlay

import "golang.org/x/sys/cpu"

//
// This file is built when we're not on arm64 or we're on windows/arm64. In the
// former case, just returning false would do. In the latter case, the right thing
// to do is to return cpu.ARM64.HasXXX. Because cpu.ARM64.HasXXX are always false
// when not on arm64, we conflate these two cases in a single file.
//

// arm64HasAES returns whether the CPU supports AES.
func arm64HasAES() bool {
	return cpu.ARM64.HasAES
}

// arm64HasPMULL returns whether the CPU supports PMULL.
func arm64HasPMULL() bool {
	return cpu.ARM64.HasPMULL
}
