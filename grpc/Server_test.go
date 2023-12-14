package grpc

import (
	"testing"
)

func TestInitGRPCServer(t *testing.T) {
	// Test with TLS enabled (mutual TLS enabled by default)
	t.Setenv("INSECURE_GRPC_SERVER", "false")
	_, err := InitGRPCServer()
	if err.Error() != "SERVER_CRT_LOCATION & SERVER_KEY_LOCATION env vars must be set" {
		t.Errorf(err.Error())
	}

	// Test with TLS disabled (mutual TLS enabled by default)
	t.Setenv("INSECURE_GRPC_SERVER", "true")
	_, err = InitGRPCServer()
	if err != nil {
		t.Errorf(err.Error())
	}
}
