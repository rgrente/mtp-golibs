package grpc

import (
	"testing"
)

func TestInitGRPCClient(t *testing.T) {
	// Test with TLS enabled
	t.Setenv("INSECURE_GRPC_CLIENT", "false")
	_, err := InitGRPCClient("localhost:10000")
	if err.Error() != "CA_CRT_LOCATION env var must be set" {
		t.Errorf(err.Error())
	}

	// Test with TLS disabled
	t.Setenv("INSECURE_GRPC_CLIENT", "true")
	_, err = InitGRPCClient("localhost:10000")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Test with TLS disabled + empty target
	t.Setenv("INSECURE_GRPC_CLIENT", "true")
	_, err = InitGRPCClient("")
	if err.Error() != "failed to build resolver: passthrough: received empty target in Build()" {
		t.Errorf(err.Error())
	}
}
