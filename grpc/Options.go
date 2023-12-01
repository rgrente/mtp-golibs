package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/rgrente/mtp-golibs/security"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func SetGRPCClientOptions() ([]grpc.DialOption, error) {

	var opts []grpc.DialOption

	// load TLS credentials (CA)
	CACertPool, err := security.SetCACertPool()
	if err != nil {
		log.Fatal("cannot generate CA cert pool: ", err)
	}

	// Create TLS config
	config := &tls.Config{
		RootCAs: CACertPool,
	}

	// append opts with TLS config credential
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))

	return opts, nil

}

func SetGRPCServerOptions(disableClientAuth bool) ([]grpc.ServerOption, error) {

	// check if SERVER_CRT_LOCATION & SERVER_KEY_LOCATION env vars are set
	if os.Getenv("SERVER_CRT_LOCATION") == "" || os.Getenv("SERVER_KEY_LOCATION") == "" {
		return nil, errors.New("SERVER_CRT_LOCATION & SERVER_KEY_LOCATION env vars must be set")
	}

	// Load server key pair
	cert, err := tls.LoadX509KeyPair(os.Getenv("SERVER_CRT_LOCATION"), os.Getenv("SERVER_KEY_LOCATION"))
	if err != nil {
		return nil, err
	}

	// Create TLS config
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// if client auth enabled
	if !disableClientAuth {

		// check if CA_CRT_LOCATION env var is set
		if os.Getenv("CA_CRT_LOCATION") == "" {
			return nil, errors.New("CA_CRT_LOCATION env var must be set when client auth is enabled")
		}

		// Load CA cert (to check client cert)
		caCert, err := os.ReadFile(os.Getenv("CA_CRT_LOCATION"))
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

	// set TLS option
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(config)),
	}

	return opts, nil

}
