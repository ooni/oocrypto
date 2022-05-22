// SPDX-License-Identifier: BSD-3-Clause

package tls

import (
	"context"
	stdlibtls "crypto/tls"
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
// This function returns a *ConnStdlib type in case of success.
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
func NewClientConnStdlib(conn net.Conn, config *stdlibtls.Config) (*ConnStdlib, error) {
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
	return &ConnStdlib{Client(conn, ourConfig)}, nil
}

// connStdlibUnderlyingConn is similar to oohttp.TLSConn but its ConnectionState
// method returns the ConnectionState defined by this package rather than the
// equivalent one defined by crypto/tls in the stdlib.
//
// We use this type instead of directly using *Conn to simplify unit
// testing of the ConnStdlib.ConnectionState method.
type connStdlibUnderlyingConn interface {
	net.Conn
	HandshakeContext(ctx context.Context) error
	ConnectionState() ConnectionState
}

// ConnStdlib is the Conn-like type returned by NewClientConnStdlib. This type
// is pretty much like this package's Conn except that the ConnectionState method
// returns crypto/tls's ConnectionState. This change is enough to make this
// struct compatible with github.com/ooni/oohttp.TLSConn.
type ConnStdlib struct {
	connStdlibUnderlyingConn
}

// connStdlibOOHTTPTLSLikeConn is equivalent to oohttp.TLSConn. We want this type to ensure
// our ConnStdlib type implements the desired interface used by OONI.
type connStdlibOOHTTPTLSLikeConn interface {
	net.Conn

	HandshakeContext(ctx context.Context) error

	ConnectionState() stdlibtls.ConnectionState
}

var _ connStdlibOOHTTPTLSLikeConn = &ConnStdlib{} // ensure we implement this interface

// ConnectionState converts the underlying Conn's ConnectionState to the
// equivalent type exported by the Go standard library.
func (c *ConnStdlib) ConnectionState() stdlibtls.ConnectionState {
	state := c.connStdlibUnderlyingConn.ConnectionState()
	return stdlibtls.ConnectionState{
		Version:                     state.Version,
		HandshakeComplete:           state.HandshakeComplete,
		DidResume:                   state.DidResume,
		CipherSuite:                 state.CipherSuite,
		NegotiatedProtocol:          state.NegotiatedProtocol,
		NegotiatedProtocolIsMutual:  state.NegotiatedProtocolIsMutual,
		ServerName:                  state.ServerName,
		PeerCertificates:            state.PeerCertificates,
		VerifiedChains:              state.VerifiedChains,
		SignedCertificateTimestamps: state.SignedCertificateTimestamps,
		OCSPResponse:                state.OCSPResponse,
		TLSUnique:                   state.TLSUnique,
	}
}
