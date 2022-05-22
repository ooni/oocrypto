// SPDX-License-Identifier: BSD-3-Clause

package tls

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"reflect"
)

// ErrIncompatibleStdlibConfig is returned when NewClientConnStdlib is
// passed an incompatible config (i.e., a config containing fields that
// we don't know how to convert and whose value is not ~zero).
var ErrIncompatibleStdlibConfig = errors.New("ootls: incompatible stdlib config")

// NewClientConnStdlib is like Client but takes in input a *crypto/tls.Config
// rather than a *github.com/ooni/oocrypto/tls.Config.
//
// The config cannot be nil: users must set either ServerName or
// InsecureSkipVerify in the config.
//
// This function will return ErrIncompatibleStdlibConfig if unsupported
// fields have a nonzero value, because the resulting Conn will not
// be compatible with the configuration you provided us with.
//
// We currently support these fields:
//
// - DynamicRecordSizingDisabled
//
// - InsecureSkipVerify
//
// - MaxVersion
//
// - MinVersion
//
// - NextProtos
//
// - RootCAs
//
// - ServerName
func NewClientConnStdlib(conn net.Conn, config *tls.Config) (*Conn, error) {
	supportedFields := map[string]bool{
		"DynamicRecordSizingDisabled": true,
		"InsecureSkipVerify":          true,
		"MaxVersion":                  true,
		"MinVersion":                  true,
		"NextProtos":                  true,
		"RootCAs":                     true,
		"ServerName":                  true,
	}
	value := reflect.ValueOf(config).Elem()
	kind := value.Type()
	for idx := 0; idx < value.NumField(); idx++ {
		field := value.Field(idx)
		if field.IsZero() {
			continue
		}
		fieldKind := kind.Field(idx)
		if supportedFields[fieldKind.Name] {
			continue
		}
		err := fmt.Errorf("%w: field %s is nonzero", ErrIncompatibleStdlibConfig, fieldKind.Name)
		return nil, err
	}
	ourConfig := &Config{
		DynamicRecordSizingDisabled: config.DynamicRecordSizingDisabled,
		InsecureSkipVerify:          config.InsecureSkipVerify,
		MaxVersion:                  config.MaxVersion,
		MinVersion:                  config.MinVersion,
		NextProtos:                  config.NextProtos,
		RootCAs:                     config.RootCAs,
		ServerName:                  config.ServerName,
	}
	return Client(conn, ourConfig), nil
}
