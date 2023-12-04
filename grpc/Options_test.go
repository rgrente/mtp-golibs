package grpc

import (
	"testing"
)

func TestSetGRPCClientOptions(t *testing.T) {

	// no en vars
	_, err := SetGRPCClientOptions()
	if err.Error() != "CA_CRT_LOCATION env var must be set" {
		t.Errorf(err.Error())
	}

	// ca wrong env var path
	t.Setenv("CA_CRT_LOCATION", "../tests/c.crt")
	_, err = SetGRPCClientOptions()
	if err.Error() != "open ../tests/c.crt: no such file or directory" {
		t.Errorf(err.Error())
	}

	// ca env ok
	t.Setenv("CA_CRT_LOCATION", "../tests/ca.crt")
	_, err = SetGRPCClientOptions()
	if err != nil {
		t.Errorf(err.Error())
	}

	// one of client crt or client key set env var path
	t.Setenv("CLIENT_KEY_LOCATION", "../tests/tls.crt")
	_, err = SetGRPCClientOptions()
	if err.Error() != "CLIENT_CRT_LOCATION and CLIENT_KEY_LOCATION must both be set to send client cert for outbound connections" {
		t.Errorf(err.Error())
	}

	// wrong client crt env var path
	t.Setenv("CLIENT_CRT_LOCATION", "../tests/tl.crt")
	t.Setenv("CLIENT_KEY_LOCATION", "../tests/tls.key")
	_, err = SetGRPCClientOptions()
	if err.Error() != "open ../tests/tl.crt: no such file or directory" {
		t.Errorf(err.Error())
	}

	// mismatch client crt & key
	t.Setenv("CLIENT_CRT_LOCATION", "../tests/ca.crt")
	t.Setenv("CLIENT_KEY_LOCATION", "../tests/tls.key")
	_, err = SetGRPCClientOptions()
	if err.Error() != "tls: private key does not match public key" {
		t.Errorf(err.Error())
	}

	// client crt & key env vars ok
	t.Setenv("CLIENT_CRT_LOCATION", "../tests/tls.crt")
	t.Setenv("CLIENT_KEY_LOCATION", "../tests/tls.key")
	_, err = SetGRPCClientOptions()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestSetGRPCServerOptions(t *testing.T) {

	// no env vars + disable client auth
	_, err := SetGRPCServerOptions(true)
	if err.Error() != "SERVER_CRT_LOCATION & SERVER_KEY_LOCATION env vars must be set" {
		t.Errorf(err.Error())
	}

	// no env vars + enable client auth
	_, err = SetGRPCServerOptions(false)
	if err.Error() != "SERVER_CRT_LOCATION & SERVER_KEY_LOCATION env vars must be set" {
		t.Errorf(err.Error())
	}

	t.Setenv("SERVER_CRT_LOCATION", "../tests/tl.crt")
	t.Setenv("SERVER_KEY_LOCATION", "../tests/tls.key")
	// wrong env vars path + disable client auth
	_, err = SetGRPCServerOptions(true)
	if err.Error() != "open ../tests/tl.crt: no such file or directory" {
		t.Errorf(err.Error())
	}

	// wrong env vars path + enable client auth
	_, err = SetGRPCServerOptions(false)
	if err.Error() != "open ../tests/tl.crt: no such file or directory" {
		t.Errorf(err.Error())
	}

	t.Setenv("SERVER_CRT_LOCATION", "../tests/tls.crt")
	t.Setenv("SERVER_KEY_LOCATION", "../tests/tls.key")
	// server env vars + disable client auth
	_, err = SetGRPCServerOptions(true)
	if err != nil {
		t.Errorf(err.Error())
	}

	// server env vars + enable client auth
	_, err = SetGRPCServerOptions(false)
	if err.Error() != "TLS_CLIENT_AUTH_CA_CRT_LOCATION env var must be set when client auth is enabled" {
		t.Errorf(err.Error())
	}

	t.Setenv("TLS_CLIENT_AUTH_CA_CRT_LOCATION", "../tests/c.crt")
	// all env vars + disable client auth
	_, err = SetGRPCServerOptions(true)
	if err != nil {
		t.Errorf(err.Error())
	}

	// all env vars + enable client auth
	_, err = SetGRPCServerOptions(false)
	if err.Error() != "open ../tests/c.crt: no such file or directory" {
		t.Errorf(err.Error())
	}

	t.Setenv("TLS_CLIENT_AUTH_CA_CRT_LOCATION", "../tests/ca.crt")
	// all env vars + disable client auth
	_, err = SetGRPCServerOptions(true)
	if err != nil {
		t.Errorf(err.Error())
	}

	// all env vars + enable client auth
	_, err = SetGRPCServerOptions(false)
	if err != nil {
		t.Errorf(err.Error())
	}
}
