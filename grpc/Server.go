package grpc

import (
	"os"

	"google.golang.org/grpc"
)

// InitGRPCServer initializes and returns a gRPC server based on options configured in environment variables.
func InitGRPCServer() (*grpc.Server, error) {
	// Set GRPC options
	var opts []grpc.ServerOption
	var err error

	// Check if the GRPC server should be configured without TLS
	if os.Getenv("INSECURE_GRPC_SERVER") == "true" {
		// TLS disabled on the GRPC server
		opts, err = SetGRPCServerOptions(true, true)
		if err != nil {
			return nil, err
		}
	} else {
		// Check if mutual TLS authentication is disabled
		if os.Getenv("TLS_CLIENT_AUTH_DISABLED") == "true" {
			// Mutual TLS authentication disabled
			opts, err = SetGRPCServerOptions(true, false)
			if err != nil {
				return nil, err
			}
		} else {
			// Mutual TLS authentication enabled (default)
			opts, err = SetGRPCServerOptions(false, false)
			if err != nil {
				return nil, err
			}
		}
	}

	// Create and return a new gRPC server with the configured options
	return grpc.NewServer(opts...), nil
}
