package security

import (
	"testing"
)

func TestLoadTLSCredentials(t *testing.T) {
	t.Setenv("CA_CRT_LOCATION", "../tests/ca.crt")
	_, err := LoadTLSCredentials()
	if err != nil {
		t.Errorf(err.Error())
	}
}
