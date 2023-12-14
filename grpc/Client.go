package grpc

import (
	"os"

	"google.golang.org/grpc"
)

// InitGRPCClient initializes and returns a gRPC client connection to the specified target.
func InitGRPCClient(target string) (*grpc.ClientConn, error) {
	// Set GRPC client options
	var opts []grpc.DialOption
	var err error

	if os.Getenv("INSECURE_GRPC_CLIENT") == "true" {
		// Insecure GRPC client (without TLS)
		opts, err = SetGRPCClientOptions(true)
		if err != nil {
			return nil, err
		}
	} else {
		// Secure GRPC client (with TLS)
		opts, err = SetGRPCClientOptions(false)
		if err != nil {
			return nil, err
		}
	}

	// Open GRPC connection
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}
	// Defer conn.Close() - Consider uncommenting this line if you want to close the connection after using it.

	return conn, nil
}
