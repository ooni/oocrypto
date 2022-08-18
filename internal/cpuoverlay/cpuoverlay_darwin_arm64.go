//go:build arm64 && darwin

package cpuoverlay

//
// As documented below, there's no support for detecting the
// capabilities of darwin/arm64 both in internal/cpu and x/sys/cpu
// so we're going to do what internal/cpu currently is doing.
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
