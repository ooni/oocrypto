// Package aes redirects to yawning/bsaes.
package aes

import "gitlab.com/yawning/bsaes.git"

// NewCipher forwards the call to [bsaes.NewCipher].
var NewCipher = bsaes.NewCipher

// BlockSize forwards the call to [bsaes.BlockSize].
var BlockSize = bsaes.BlockSize
