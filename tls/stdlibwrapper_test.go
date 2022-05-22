// SPDX-License-Identifier: BSD-3-Clause

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"testing"
	"time"
)

func TestNewClientConnStdlib(t *testing.T) {
	tests := []struct {
		name   string
		config *tls.Config
		err    error
	}{{
		name: "with only supported config fields",
		config: &tls.Config{
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
		config: &tls.Config{
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
