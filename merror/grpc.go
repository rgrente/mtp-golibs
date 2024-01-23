package merror

import (
	"context"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ProcessGRPCError(ctx context.Context, e error) error {
	var err *MError
	var trailer metadata.MD
	var trace string
	serviceName, ok := os.LookupEnv("SERVICE_NAME")
	pc, file, line, callerOk := runtime.Caller(1)

	switch {
	case errors.As(e, &err):
		if !ok {
			if err.Trace == "" && callerOk {
				trace = file[strings.LastIndex(file, "/")+1:] +
				":" + runtime.FuncForPC(pc).Name() +
				":" + strconv.Itoa(line)
			} else {
				trace = err.Trace
			}
		} else {
			trace = serviceName+"/"+err.Trace
		}
	default:
		err = err.New(e.Error())
		if !ok {
			trace = ""
		} else {
			trace = serviceName
		}
	}
	trailer = metadata.Pairs(
		"code", strconv.Itoa(err.Code),
		"message", err.Message,
		"description", err.Description,
		"trace", trace,
		"dev_code", strconv.Itoa(err.DevCode),
		"dev_message", err.DevMessage,
		"dev_description", err.DevDescription,
	)

	grpc.SetTrailer(ctx, trailer)
	return status.Error(codes.Aborted, err.Message)
}
