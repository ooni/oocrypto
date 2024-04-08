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

You SHOULD use this package with the exact Go version from which we extracted
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

We changed approach with [#26](https://github.com/ooni/oocrypto/pull/26) to
use [gitlab.com/yawning/bsaes](https://gitlab.com/yawning/bsaes) for constant
time AES.

## Update procedure

(Adapted from ooni/oohttp instructions.)

- [ ] update [UPSTREAM](UPSTREAM), commit the change, and then
run the `./tools/merge.bash` script to merge from upstream;

- [ ] fix all the likely merge conflicts

- [ ] delete every package except for `tls`

- [ ] ensure that `stdlibwrapper.go` correctly fills `tls.ConnectionState`
in the `ConnStdlib.ConnectionState` method

- [ ] use `./tools/compare.bash` to make sure the changes with respect
to upstream are reasonable

- [ ] `go build -v ./...` must succeed

- [ ] `go test -race ./...` must succeed

- [ ] run `go get -u -v ./... && go mod tidy`

- [ ] open a pull request using this check-list as its content and merge it preserving history
