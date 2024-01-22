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
	me "github.com/rgrente/mtp-golibs/merror"
)

func processGRPCError(ctx context.Context, e error) error {
	var err *me.MError
	var trailer metadata.MD
	var trace string
	serviceName, ok := os.LookupEnv("SERVICE_NAME")
	pc, file, line, callerOk := runtime.Caller(1)

	switch {
	case errors.As(e, &err):
		if !ok {
			if err.trace == "" && callerOk {
				trace = file[strings.LastIndex(file, "/")+1:] +
				":" + runtime.FuncForPC(pc).Name() +
				":" + strconv.Itoa(line)
			} else {
				trace = err.trace
			}
		} else {
			trace = serviceName+"/"+err.trace
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
		"code", strconv.Itoa(err.code),
		"message", err.message,
		"description", err.description,
		"trace", trace,
		"dev_code", strconv.Itoa(err.devCode),
		"dev_message", err.devMessage,
		"dev_description", err.devDescription,
	)

	grpc.SetTrailer(ctx, trailer)
	return status.Error(codes.Aborted, err.message)
}
