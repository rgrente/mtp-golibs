package merror

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

type MError struct {
	Code int
	Message string
	Description string
	DevCode int
	DevMessage string
	DevDescription string
	Trace string
}

func (e *MError) Error() string {
	return e.Message
}

func (e *MError) New(message string) *MError {
	var err = MError{}
	var callerOk bool
	pc, file, line, callerOk := runtime.Caller(1)

	err.Code = 501
	err.Message = "Unknown error"
	err.Description = ""
	err.DevCode = 1000
	err.DevMessage = ""
	err.DevDescription = ""
	if callerOk {
		err.Trace = file[strings.LastIndex(file, "/")+1:] +
		":" + runtime.FuncForPC(pc).Name() +
		":" + strconv.Itoa(line)
	} else {
		err.Trace = ""
	}
	return &err
}

func ProcessError(e error) error {
	var err *MError
	var callerOk bool
	pc, file, line, callerOk := runtime.Caller(1)

	switch {
	case errors.As(e, &err):
		if err.Trace == "" {
			err.Trace = file[strings.LastIndex(file, "/")+1:] +
			":" + runtime.FuncForPC(pc).Name() +
			":" + strconv.Itoa(line)
		}
		if callerOk {
			err.Trace = runtime.FuncForPC(pc).Name() + "/" + err.Trace
		}
	default:
		err = err.New(e.Error())
	}
	return err
}
