// Package cpuoverlay is a logic overlay on top of x/sys/cpu.
//
// The main problem that we want to solve is that on Android there are
// cases where reading /proc/self/auxv is not possible.
//
// This causes crypto/tls to not choose AES where it would otherwise
// be possible, in turn causing censorship. See also the
// https://github.com/ooni/probe/issues/1444 issue for more details.
//
// Ideally we would like to call getauxval(3) when initializing
// the runtime package. However, runtime cannot use CGO. Doing that
// leads to an import loop, so we cannot build.
//
// We could also try to parse /proc/cpuinfo (I didn't explore this route).
//
// The solution chosen here is to export predicates on the CPU
// functionality. We limit ourselves to what we need in order to
// choose AES in crypto/tls when the CPU supports it.
//
// Additionally, we may use this package to solve other CPU issues.
//
// This package defines GOOS/GOARCH specific files with predicate
// functions for the architectures for which we need to overlay over
// the functionality provided by the x/sys/cpu package.
package cpuoverlay

// Arm64HasAES returns whether the CPU supports AES.
func Arm64HasAES() bool {
	return arm64HasAES()
}

// Arm64HasPMULL returns whether the CPU supports PMULL.
func Arm64HasPMULL() bool {
	return arm64HasPMULL()
}
