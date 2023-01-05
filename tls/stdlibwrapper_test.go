// SPDX-License-Identifier: BSD-3-Clause

package tls

import (
	"context"
	stdlibtls "crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestNewClientConnStdlib(t *testing.T) {
	tests := []struct {
		name   string
		config *stdlibtls.Config
		err    error
	}{{
		name: "with only supported config fields",
		config: &stdlibtls.Config{
			DynamicRecordSizingDisabled: true,
			RootCAs:                     x509.NewCertPool(),
			ServerName:                  "ooni.org",
			InsecureSkipVerify:          true,
			MinVersion:                  VersionTLS10,
			MaxVersion:                  VersionTLS13,
			NextProtos:                  []string{"h3"},
		},
		err: nil,
	}, {
		name: "with unsupported fields",
		config: &stdlibtls.Config{
			Time: func() time.Time {
				return time.Now()
			},
		},
		err: ErrIncompatibleStdlibConfig,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := net.Dial("udp", "8.8.8.8:443") // we just want a valid conn
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			got, err := NewClientConnStdlib(conn, tt.config)
			if !errors.Is(err, tt.err) {
				t.Fatal("unexpected error", err)
			}
			if err == nil && got == nil {
				t.Fatal("expected non-nil conn here")
			}
		})
	}
}

// connStdlibUnderlyingConnMockable allows to mock connStdlibUnderlyingConn.
type connStdlibUnderlyingConnMockable struct {
	// Conn is the embedded mockable Conn.
	net.Conn

	// MockConnectionState allows to mock the ConnectionState method.
	MockConnectionState func() ConnectionState

	// MockHandshakeContext allows to mock the HandshakeContext method.
	MockHandshakeContext func(ctx context.Context) error

	// MockNetConn allows to mock the NetConn method
	MockNetConn func() net.Conn
}

// ConnectionState calls MockConnectionState.
func (c *connStdlibUnderlyingConnMockable) ConnectionState() ConnectionState {
	return c.MockConnectionState()
}

// HandshakeContext calls MockHandshakeContext.
func (c *connStdlibUnderlyingConnMockable) HandshakeContext(ctx context.Context) error {
	return c.MockHandshakeContext(ctx)
}

// NetConn calls MockNetConn
func (c *connStdlibUnderlyingConnMockable) NetConn() net.Conn {
	return c.MockNetConn()
}

func TestConnStdlib_ConnectionState(t *testing.T) {
	type fields struct {
		connLike connStdlibUnderlyingConn
	}
	tests := []struct {
		name   string
		fields fields
		want   stdlibtls.ConnectionState
	}{{
		name: "common case",
		fields: fields{
			connLike: &connStdlibUnderlyingConnMockable{
				MockConnectionState: func() ConnectionState {
					return ConnectionState{
						Version:                    1,
						HandshakeComplete:          true,
						DidResume:                  true,
						CipherSuite:                2,
						NegotiatedProtocol:         "echo",
						NegotiatedProtocolIsMutual: true,
						ServerName:                 "x.org",
						PeerCertificates: []*x509.Certificate{{
							Raw: []byte{0x01},
						}},
						VerifiedChains: [][]*x509.Certificate{{{
							Raw: []byte{0x77},
						}}},
						SignedCertificateTimestamps: [][]byte{{0x44}},
						OCSPResponse:                []byte{0x14, 0x11},
						TLSUnique:                   []byte{0x17},
					}
				},
			},
		},
		want: stdlibtls.ConnectionState{
			Version:                    1,
			HandshakeComplete:          true,
			DidResume:                  true,
			CipherSuite:                2,
			NegotiatedProtocol:         "echo",
			NegotiatedProtocolIsMutual: true,
			ServerName:                 "x.org",
			PeerCertificates: []*x509.Certificate{{
				Raw: []byte{0x01},
			}},
			VerifiedChains: [][]*x509.Certificate{{{
				Raw: []byte{0x77},
			}}},
			SignedCertificateTimestamps: [][]byte{{0x44}},
			OCSPResponse:                []byte{0x14, 0x11},
			TLSUnique:                   []byte{0x17},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConnStdlib{
				connStdlibUnderlyingConn: tt.fields.connLike,
			}
			if got := c.ConnectionState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnStdlib.ConnectionState() = %v, want %v", got, tt.want)
			}
		})
	}
}
