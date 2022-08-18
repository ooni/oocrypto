// Package cpuoverlay is a logic overlay on top of x/sys/cpu that
// attempts to avoid cases in which x/sys/cpu is wrong.
//
// android/arm64
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
// Until this is fixed in src/runtime, we call getauxval(3) here.
//
// darwin/arm64
//
// Additionally, we may use this package to solve other CPU issues. For
// example, x/sys/cpu does not currently know about darwin/arm64 features.
//
// Design
//
// This package contains GOOS/GOARCH-specific files with predicate
// functions returning the correct value in cases in which x/sys/cpu
// is wrong. You are expected to replace the code that normally
// lives inside src/crypto/tls and src/cryto/aes and that we have
// forked in this repository such that you call the predicates
// of this package as opposed to using directly x/sys/cpu values.
package cpuoverlay

// We dispatch to GOOS/GOARCH specific implementations

// Arm64HasAES returns whether the CPU supports AES.
func Arm64HasAES() bool {
	return arm64HasAES()
}

// Arm64HasPMULL returns whether the CPU supports PMULL.
func Arm64HasPMULL() bool {
	return arm64HasPMULL()
}
