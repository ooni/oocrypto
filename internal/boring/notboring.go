// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !boringcrypto || !linux || !amd64 || !cgo || android || cmd_go_bootstrap || msan
// +build !boringcrypto !linux !amd64 !cgo android cmd_go_bootstrap msan

package boring

import "crypto/cipher"

const available = false

// Unreachable marks code that should be unreachable
// when BoringCrypto is in use. It is a no-op without BoringCrypto.
func Unreachable() {
	// nothing
}

func NewGCMTLS(cipher.Block) (cipher.AEAD, error) { panic("boringcrypto: not available") }
