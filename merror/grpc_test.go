package merror

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProcessGRPCError(t *testing.T) {
	var err = &MError{}
	err.Code = 400
	err.Message = "Bad Request"
	err.Description = "Invalid argument"
	err.DevCode = 1001
	err.DevMessage = "Bad Request"
	err.DevDescription = "Invalid argument"
	ctx := context.Background()
	e := ProcessGRPCError(ctx, err)
	if e, ok := status.FromError(e); ok {
		if e.Code() != codes.Aborted {
			t.Errorf("Output Code %s not equal to expected %s", e.Code().String(), codes.Aborted.String())
		} else if e.Message() != err.Message {
			t.Errorf("Output Message %s not equal to expected %s", e.Message(), err.Message)
		}
	}
}