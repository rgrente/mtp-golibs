package merror

import (
	"errors"
	"testing"
	"net/http/httptest"

	"google.golang.org/grpc/metadata"
	"github.com/gin-gonic/gin"
)

func TestMError(t *testing.T) {
	trailer := metadata.Pairs(
		"code", "400",
		"message", "Bad Request",
		"description", "Invalid argument",
		"trace", "",
		"dev_code", "1001",
		"dev_message", "Bad Request",
		"dev_description", "Invalid argument",
	)
	err := ToMError(trailer)
	merr := &MError{}
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
		}
	} else {
		t.Errorf("Returned error from Process error is not of type MError")
	}
}

func TestRenderError(t *testing.T) {
	err := errors.New("")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	RenderError(c, err)
}