package grpc

import (
	"crypto/tls"
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

	// Load certs
	cert, err := tls.LoadX509KeyPair(os.Getenv("CLIENT_RS_SERVER_CRT_LOCATION"), os.Getenv("CLIENT_RS_SERVER_KEY_LOCATION"))
	if err != nil {
		return nil, err
	}

	// set TLS option
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	return opts, nil

}
