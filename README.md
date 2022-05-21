# github.com/ooni/oocrypto

This repository contains a fork of the Go standard library's `crypto`
package including patches to allow selecting AES hardware support
on Android devices. We documented why we need these patches at OONI in
the [Making the OONI Probe Android app more resilient](
https://ooni.org/post/making-ooni-probe-android-more-resilient/) blog post.

## Motivation and maintenance

To solve our issues with Android apps, we originally forked golang/go
itself at [ooni/go](https://github.com/ooni/go). However, a full fork of
Go required us to compile this fork and build Android apps using it,
which was making building OONI excessively complicated. Hence, we later
chose to just fork the `crypto` package. We will continue to keep this
fork up to date as long as it serves our goals.

## License

Each individual file from the `crypto` fork maintains its original
copyright and any change to such files authored by us keeps the same
BSD license of the original code. Because we anticipate integrating
code under the GPL license from `Yawning/utls` we chose to license the
repository using the GPL.

```
SPDX-License-Identifier: GPL-3.0-or-later
```

## Issue tracker

Please, report issues in the [ooni/probe](https://github.com/ooni/probe)
repository. Make sure you mention `oocrypto` in the issue title.

## Update procedure

(Adapted from ooni/oohttp instructions.)

- [ ] run the following commands:

```bash
set -ex
git checkout main
git remote add golang git@github.com:golang/go.git || git fetch golang
git branch -D golang-upstream golang-http-upstream merged-main || true
git fetch golang
git checkout -b golang-upstream go1.17.10
git subtree split -P src/crypto/ -b golang-http-upstream
git checkout main
git checkout -b merged-main
git merge golang-http-upstream
```

- [ ] fix all the likely merge conflicts

- [ ] delete all the new packages we can safely delete. We can safely
delete a package if the package is not `tls` and:

1. either the package does not depend on `internal/cpu`

2. or the documentation of the package does not explicitly state that
the package is only secure depending on the CPU configuration, which
currently only holds for `aes` (see [aes/const.go](aes/const.go))

- [ ] `go build -v ./...` must succeed

- [ ] `go test -race ./...` must succeed

- [ ] open a pull request and merge it preserving history
