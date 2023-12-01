package grpc

import "testing"

func TestSetGRPCClientOptions(t *testing.T) {
	t.Setenv("CA_CRT_LOCATION", "../tests/ca.crt")
	_, err := SetGRPCClientOptions()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestSetGRPCServerOptions(t *testing.T) {
	t.Setenv("SERVER_CRT_LOCATION", "../tests/tls.crt")
	t.Setenv("SERVER_KEY_LOCATION", "../tests/tls.key")
	_, err := SetGRPCServerOptions()
	if err != nil {
		t.Errorf(err.Error())
	}
}
