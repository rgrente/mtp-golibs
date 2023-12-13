package grpc

import (
	"os"

	"google.golang.org/grpc"
)

func InitGRPCServer(listenPort int) (*grpc.Server, error) {

	// set GRPC options
	var opts []grpc.ServerOption
	var err error
	if os.Getenv("INSECURE_GRPC_SERVER") == "true" {
		// TLS on GRPC server disabled
		opts, err = SetGRPCServerOptions(true, true)
		if err != nil {
			return nil, err
		}
	} else {
		if os.Getenv("TLS_CLIENT_AUTH_DISABLED") == "true" {
			// mutual TLS disabled
			opts, err = SetGRPCServerOptions(true, false)
			if err != nil {
				return nil, err
			}
		} else {
			// mutual TLS enabled (default)
			opts, err = SetGRPCServerOptions(false, false)
			if err != nil {
				return nil, err
			}
		}
	}

	return grpc.NewServer(opts...), nil
}
