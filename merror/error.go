package merror

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

// MError represents a custom error type with detailed information for developers.
type MError struct {
	Code             int    // Code represents the error code.
	Message          string // Message represents a user-friendly error message.
	Description      string // Description provides additional details about the error.
	DevCode          int    // DevCode represents the developer-specific error code.
	DevMessage       string // DevMessage is a message intended for developers.
	DevDescription   string // DevDescription is additional information for developers.
	Trace            string // Trace provides information about the source of the error.
}

// Error returns the user-friendly error message when the MError type is used as an error.
func (e *MError) Error() string {
	return e.Message
}

// New creates a new MError instance with default values and sets the error message.
// It also captures the file, function, and line number where the error occurred.
func (e *MError) New(message string) *MError {
	var err = MError{}
	var callerOk bool
	pc, file, line, callerOk := runtime.Caller(1)

	// Set default values for the MError instance.
	err.Code = 501
	err.Message = "Unknown error"
	err.Description = ""
	err.DevCode = 1000
	err.DevMessage = ""
	err.DevDescription = ""

	// If runtime.Caller was successful, capture file, function, and line number in Trace.
	if callerOk {
		err.Trace = file[strings.LastIndex(file, "/")+1:] +
			":" + runtime.FuncForPC(pc).Name() +
			":" + strconv.Itoa(line)
	} else {
		err.Trace = ""
	}

	return &err
}

// ProcessError processes an error and returns an MError with additional developer-specific information.
// If the input error is not an MError, it creates a new MError with the error message.
func ProcessError(e error) error {
	var err *MError
	var callerOk bool
	pc, file, line, callerOk := runtime.Caller(1)

	switch {
	case errors.As(e, &err):
		// If the input error is already an MError, use it and update the Trace if necessary.
		if err.Trace == "" {
			err.Trace = file[strings.LastIndex(file, "/")+1:] +
				":" + runtime.FuncForPC(pc).Name() +
				":" + strconv.Itoa(line)
		}
		if callerOk {
			err.Trace = runtime.FuncForPC(pc).Name() + "/" + err.Trace
		}
	default:
		// If the input error is not an MError, create a new MError with the error message.
		err = err.New(e.Error())
	}

	return err
}
