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

// ProcessGRPCError processes a gRPC error, enriches it with additional information, and returns a gRPC status error.
// Each time a function gets an error: using this function before return will add trace of the current function in the error
func ProcessGRPCError(ctx context.Context, e error) error {
	var err *MError
	var trailer metadata.MD
	var trace string
	serviceName, ok := os.LookupEnv("SERVICE_NAME")
	pc, file, line, callerOk := runtime.Caller(1)

	// Check if the input error is of type MError.
	switch {
	case errors.As(e, &err):
		// If SERVICE_NAME is not set, use either the caller information or the trace from the MError.
		if !ok {
			if err.Trace == "" && callerOk {
				trace = file[strings.LastIndex(file, "/")+1:] +
					":" + runtime.FuncForPC(pc).Name() +
					":" + strconv.Itoa(line)
			} else {
				trace = err.Trace
			}
		} else {
			// If SERVICE_NAME is set, use it as a prefix for the trace.
			trace = serviceName + "/" + err.Trace
		}
	default:
		// If the input error is not of type MError, create a new MError with the error message.
		err = err.New(e.Error())
		// If SERVICE_NAME is not set, the trace is empty; otherwise, use SERVICE_NAME as the trace.
		if !ok {
			trace = ""
		} else {
			trace = serviceName
		}
	}

	// Create gRPC metadata trailer with MError details.
	trailer = metadata.Pairs(
		"code", strconv.Itoa(err.Code),
		"message", err.Message,
		"description", err.Description,
		"trace", trace,
		"dev_code", strconv.Itoa(err.DevCode),
		"dev_message", err.DevMessage,
		"dev_description", err.DevDescription,
	)

	// Set the trailer in the gRPC context and return a gRPC status error.
	grpc.SetTrailer(ctx, trailer)
	return status.Error(codes.Aborted, err.Message)
}
