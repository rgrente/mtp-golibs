package security

import (
	"crypto/x509"
	"fmt"
	"os"
)

// SetCACertPool loads a CA certificate from the specified location and creates a certificate pool.
func SetCACertPool() (*x509.CertPool, error) {
	// Load CA cert
	caCert, err := os.ReadFile(os.Getenv("CA_CRT_LOCATION"))
	if err != nil {
		return nil, err
	}

	// Generate cert pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA's certificate")
	}

	return certPool, nil
}
