package api

import (
	"crypto/tls"
	"fmt"

	_ "embed"
)

// We use these functions to supply TLS for various services that require it. To make development easy
// we bake in general localhost certs for quick bootstrap. The server will not start with dev certs loaded
// unless explicitly told to do so with devmode=true.

//go:embed localhost.crt
var devtlscert []byte

//go:embed localhost.key
var devtlskey []byte

// generateTLSConfig returns TLS config object necessary for HTTPS loaded from files. If server is in devmode and
// no cert is provided it instead loads certificates from embedded files for ease of development.
func (api *API) generateTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	var serverCert tls.Certificate
	var err error

	if api.config.Development.UseLocalhostTLS {
		serverCert, err = tls.X509KeyPair(devtlscert, devtlskey)
		if err != nil {
			return nil, err
		}
	} else {
		if certPath == "" || keyPath == "" {
			return nil, fmt.Errorf("TLS cert and key cannot be empty")
		}

		serverCert, err = tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, err
		}
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return tlsConfig, nil
}
