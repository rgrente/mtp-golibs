package security

import (
	"testing"
)

func TestSetCACertPool(t *testing.T) {
	t.Setenv("CA_CRT_LOCATION", "../tests/ca.crt")
	_, err := SetCACertPool()
	if err != nil {
		t.Errorf(err.Error())
	}
}
