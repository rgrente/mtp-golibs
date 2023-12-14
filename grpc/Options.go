package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/rgrente/mtp-golibs/security"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// SetGRPCClientOptions configures gRPC client options based on the specified parameters.
func SetGRPCClientOptions(insecureClient bool) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	if insecureClient {
		// Insecure client (without TLS)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Check if CA_CRT_LOCATION env var is set
		if os.Getenv("CA_CRT_LOCATION") == "" {
			return nil, errors.New("CA_CRT_LOCATION env var must be set")
		}

		// Load TLS credentials (CA)
		CACertPool, err := security.SetCACertPool()
		if err != nil {
			return nil, err
		}

		// Create TLS config
		config := &tls.Config{
			RootCAs: CACertPool,
		}

		// Load client key pair (for client auth) if env vars provided
		if os.Getenv("CLIENT_CRT_LOCATION") != "" && os.Getenv("CLIENT_KEY_LOCATION") != "" {
			clientCert, err := tls.LoadX509KeyPair(os.Getenv("CLIENT_CRT_LOCATION"), os.Getenv("CLIENT_KEY_LOCATION"))
			if err != nil {
				return nil, err
			}
			config.Certificates = []tls.Certificate{clientCert}
		} else if os.Getenv("CLIENT_CRT_LOCATION") != "" || os.Getenv("CLIENT_KEY_LOCATION") != "" {
			return nil, errors.New("CLIENT_CRT_LOCATION and CLIENT_KEY_LOCATION must both be set to send client cert for outbound connections")
		}

		// Append opts with TLS config credential
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}

	return opts, nil
}

// SetGRPCServerOptions configures gRPC server options based on the specified parameters.
func SetGRPCServerOptions(disableClientAuth bool, insecureServer bool) ([]grpc.ServerOption, error) {
	var opts []grpc.ServerOption

	if !insecureServer {
		// Check if SERVER_CRT_LOCATION & SERVER_KEY_LOCATION env vars are set
		if os.Getenv("SERVER_CRT_LOCATION") == "" || os.Getenv("SERVER_KEY_LOCATION") == "" {
			return nil, errors.New("SERVER_CRT_LOCATION & SERVER_KEY_LOCATION env vars must be set")
		}

		// Load server key pair
		serverCert, err := tls.LoadX509KeyPair(os.Getenv("SERVER_CRT_LOCATION"), os.Getenv("SERVER_KEY_LOCATION"))
		if err != nil {
			return nil, err
		}

		// Create TLS config
		config := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
		}

		// If client auth enabled
		if !disableClientAuth {
			// Check if TLS_CLIENT_AUTH_CA_CRT_LOCATION env var is set
			if os.Getenv("TLS_CLIENT_AUTH_CA_CRT_LOCATION") == "" {
				return nil, errors.New("TLS_CLIENT_AUTH_CA_CRT_LOCATION env var must be set when client auth is enabled")
			}

			// Load CA cert (to check client cert)
			caCert, err := os.ReadFile(os.Getenv("TLS_CLIENT_AUTH_CA_CRT_LOCATION"))
			if err != nil {
				return nil, err
			}
			certPool := x509.NewCertPool()
			if !certPool.AppendCertsFromPEM(caCert) {
				return nil, fmt.Errorf("failed to add CA's certificate")
			}

			config.ClientCAs = certPool
			config.ClientAuth = tls.RequireAndVerifyClientCert
		}

		opts = append(opts, grpc.Creds(credentials.NewTLS(config)))
	}

	return opts, nil
}
