package grpc

import (
	"crypto/tls"
	"crypto/x509"
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
	tlsCredentials, err := security.LoadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))

	return opts, nil

}

func SetGRPCServerOptions() ([]grpc.ServerOption, error) {

	// Load server key pair
	cert, err := tls.LoadX509KeyPair(os.Getenv("SERVER_CRT_LOCATION"), os.Getenv("SERVER_KEY_LOCATION"))
	if err != nil {
		return nil, err
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

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	// set TLS option
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(config)),
	}

	return opts, nil

}
