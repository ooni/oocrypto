# github.com/ooni/oocrypto

This repository contains a fork of a subset of the Go stdlib's `crypto`
package including patches to allow selecting AES hardware support
on Android devices. We documented why we need these patches at OONI in
the [Making the OONI Probe Android app more resilient](
https://ooni.org/post/making-ooni-probe-android-more-resilient/) blog post.

## Motivation and maintenance

To solve our issues with Android apps, we originally forked golang/go
itself at [ooni/go](https://github.com/ooni/go). However, a full fork of
Go required us to compile this fork and build Android apps using it,
which was making building OONI excessively complicated. Hence, we later
chose to just fork the `crypto` package and documented our efforts at
[ooni/probe#2106](https://github.com/ooni/probe/issues/2106). We
will continue to keep this fork up to date as long as it serves our goals.

## Intended usage

You MUST use this package with the exact Go version from which we extracted
the source, which is documented in the [Update procedure](#update-procedure) section. The
standard library is composed of tightly integrated packages, hence
using this code with another Go version could cause subtle security issues.

The [tls/stdlibwrapper.go](tls/stdlibwrapper.go) file contains an API that allows
converting code using `crypto/tls` to code using this package.

```Go
func NewClientConnStdlib(conn net.Conn, config *stdlibtls.Config) (*ConnStdlib, error)
```

The `NewClientConnStdlib` creates a new client conn taking in input a
`tls.Config` struct as exposed by the stdlib `crypto/tls` package. The
function returns error if you passed in config fields that we don't
know (yet?) how to convert from their stdlib definition to the equivalent
definition of `Config` implemented by this module.

The returned `ConnStdlib` type implements the following interface, which
is equivalent to [oohttp](https://github.com/ooni/oohttp)'s `TLSConn`:

```Go
import (
    "context"
    "crypto/tls"
)

type TLSConn interface {
    net.Conn

    HandshakeContext(ctx context.Context) error

    ConnectionState() tls.ConnectionState

    NetConn() net.Conn
}
```

These changes are sufficient for OONI to use this library instead
of using `crypto/tls` as the underlying TLS library.

## License

Each individual file from the `crypto` fork maintains its original
copyright and [license](https://github.com/golang/go/blob/master/LICENSE). Any
change to such files authored by us keeps the same 3-clause BSD license of
the original code. Because we anticipate integrating code under the GPL license
from `Yawning/utls` we chose to license the repository using the GPL.

```
SPDX-License-Identifier: GPL-3.0-or-later
```

## Issue tracker

Please, report issues in the [ooni/probe](https://github.com/ooni/probe)
repository. Make sure you mention `oocrypto` in the issue title.

## Patches

Commit [1137f34](https://github.com/ooni/oocrypto/commit/1137f34fc78f7b5165a37f290e0b1c5e2fb074ac)
merged go1.17.10 `src/crypto`'s subtree into this repository.

[Subsequent commits](https://github.com/ooni/oocrypto/compare/1137f34fc78f7b5165a37f290e0b1c5e2fb074ac...f09fe46bcb80d2e747b0c0ea9a2835e70710690c)
removed unused code and established a procedure to sync with upstream. As part
of these commits, we replaced `internal/cpu` with `golang.org/x/sys/cpu`.

Finally, we landed [patches](https://github.com/ooni/oocrypto/compare/f09fe46bcb80d2e747b0c0ea9a2835e70710690c...4dff9e0864cd49113a36ac8112cf887cbe215d54)
to improve hardware capability detection on `android/arm64`.

## Update procedure

(Adapted from ooni/oohttp instructions.)

- [ ] check whether hardware capability detection has been improved upstream
by reading [os_linux.go](https://github.com/golang/go/blob/go1.19.4/src/runtime/os_linux.go#L238)
and update the link to `os_linux.go` based on the upstream version that
we're tracking with this fork

- [ ] update [UPSTREAM](UPSTREAM), commit the change, and then
run the `./tools/merge.bash` script to merge from upstream;

- [ ] fix all the likely merge conflicts

- [ ] delete all the new packages we can safely delete. We can safely
delete a package if the package is not `tls` and:

1. either the package does not depend on `internal/cpu`

2. or the documentation of the package does not explicitly state that
the package is only secure depending on the CPU configuration, which
currently only holds for `aes` (see [aes/const.go](aes/const.go))

- [ ] ensure that every forked package is never imported by using
the following checks (we could also use `go list` as follows
`GOOS=os GOARCH=arch go list --json ./...`):

1. `git grep 'subtle"'`

2. `git grep 'tls"'`

3. `git grep 'aes"'`

- [ ] double check whether we need to add more checks to the list above (you
can get a list of packages using `tree -d`)

- [ ] ensure that `stdlibwrapper.go` correctly fills `tls.ConnectionState`
in the `ConnStdlib.ConnectionState` method

- [ ] `go build -v ./...` must succeed

- [ ] `go test -race ./...` must succeed

- [ ] open a pull request and merge it preserving history
