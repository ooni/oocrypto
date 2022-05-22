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
chose to just fork the `crypto` package. We will continue to keep this
fork up to date as long as it serves our goals.

## Intended usage

You MUST use this package with the exact Go version from which we extracted
the source. The standard library is composed of tightly integrated packages, hence
using this code with another Go version could cause subtle security issues.

You can find the version to which this commit is bound by checking the
update script listed in the [Update procedure](#update-procedure) section.

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
removed unused code and established a procedure to sync with upstream.

Finally, we landed [patches](https://github.com/ooni/oocrypto/compare/f09fe46bcb80d2e747b0c0ea9a2835e70710690c...4dff9e0864cd49113a36ac8112cf887cbe215d54)
to improve hardware capability detection on `android/arm64`.

## Update procedure

(Adapted from ooni/oohttp instructions.)

- [ ] check whether hardware capability detection has been improved upstream

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

- [ ] ensure that every forked package is never imported by using
the following checks (we could also use `go list`, in principle, but
I have not find a way for getting results for all architectures)

1. `git grep 'subtle"'`

2. `git grep 'tls"'`

3. `git grep 'aes"'`

- [ ] double check whether we need to add more checks to the list above

- [ ] `go build -v ./...` must succeed

- [ ] `go test -race ./...` must succeed

- [ ] open a pull request and merge it preserving history
