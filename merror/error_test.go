package merror

import (
	"errors"
	"testing"
)

func TestProcessError(t *testing.T) {
	err := errors.New("test")
	err = ProcessError(err)
	var merr *MError

	if errors.As(err, &merr) {
		if merr.Code != 501 {
			t.Errorf("Output Code %d not equal to expected %d", merr.Code, 501)
		} else if merr.Message != "Unknown error" {
			t.Errorf("Output Message %s not equal to expected %s", merr.Message, "Unknown error")
		} else if merr.Description != "" {
			t.Errorf("Output Description %s not equal to expected %s", merr.Description, "")
		} else if merr.DevCode != 1000 {
			t.Errorf("Output DevCode %d not equal to expected %d", merr.DevCode, 1000)
		} else if merr.DevMessage != "" {
			t.Errorf("Output DevMessage %s not equal to expected %s", merr.DevMessage, "")
		} else if merr.DevDescription != "test" {
			t.Errorf("Output DevDescription %s not equal to expected %s", merr.DevDescription, "test")
		}
	} else {
		t.Errorf("Returned error from Process error is not of type MError")
	}
	merr.Code = 400
	merr.Message = "Bad Request"
	merr.Description = "Invalid argument"
	merr.DevCode = 1001
	merr.DevMessage = "Bad Request"
	merr.DevDescription = "Invalid argument"
	err = ProcessError(merr)
	merr = &MError{}
	if errors.As(err, &merr) {
		if merr.Code != 400 {
			t.Errorf("Output Code %d not equal to expected %d", merr.Code, 400)
		} else if merr.Message != "Bad Request" {
			t.Errorf("Output Message %s not equal to expected %s", merr.Message, "Bad Request")
		} else if merr.Description != "Invalid argument" {
			t.Errorf("Output Description %s not equal to expected %s", merr.Description, "Invalid argument")
		} else if merr.DevCode != 1001 {
			t.Errorf("Output DevCode %d not equal to expected %d", merr.DevCode, 1001)
		} else if merr.DevMessage != "Bad Request" {
			t.Errorf("Output DevMessage %s not equal to expected %s", merr.DevMessage, "Bad Request")
		} else if merr.DevDescription != "Invalid argument" {
			t.Errorf("Output DevDescription %s not equal to expected %s", merr.DevDescription, "Invalid argument")
		} else if merr.Trace != "github.com/rgrente/mtp-golibs/merror.TestProcessError/error.go:github.com/rgrente/mtp-golibs/merror.ProcessError:73" {
			t.Errorf("Output DevDescription %s not equal to expected %s", merr.Trace, "github.com/rgrente/mtp-golibs/merror.TestProcessError/error.go:github.com/rgrente/mtp-golibs/merror.ProcessError:73")
		}
	} else {
		t.Errorf("Returned error from Process error is not of type MError")
	}
}
