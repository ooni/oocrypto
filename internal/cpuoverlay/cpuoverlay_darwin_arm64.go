//go:build arm64 && darwin

package cpuoverlay

//
// See https://github.com/ooni/probe/issues/2122
//
// See https://github.com/golang/go/issues/43046
//

// arm64HasAES returns whether the CPU supports AES.
func arm64HasAES() bool {
	return true
}

// arm64HasPMULL returns whether the CPU supports PMULL.
func arm64HasPMULL() bool {
	return true
}
