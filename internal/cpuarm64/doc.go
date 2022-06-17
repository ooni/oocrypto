// Package cpuarm64 is a getauxval aware proxy for arm64 CPUs.
//
// The problem that we want to solve is that on Android there are
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
// This package is only a replacement for arm64. We use x/sys/cpu for
// all arm64 systems but Android where we call getauxval(3), *except* for
// the concrete case of darwin/arm64, since the x/sys/cpu support for
// darmin/arm64 is not implemented:
// https://github.com/golang/go/issues/43046

package cpuarm64
