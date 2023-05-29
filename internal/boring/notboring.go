// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !(boringcrypto && linux && (amd64 || arm64) && !android && !cmd_go_bootstrap && !msan && cgo)

package boring

const available = false

// Unreachable marks code that should be unreachable
// when BoringCrypto is in use. It is a no-op without BoringCrypto.
func Unreachable() {
	// nothing
}
