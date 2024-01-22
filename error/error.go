package merror

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)


type MError struct {
	code int
	message string
	description string
	devCode int
	devMessage string
	devDescription string
	trace string
}

func (e *MError) Error() string {
	return e.message
}

func (e *MError) New(message string) *MError {
	var err *MError
	var callerOk bool
	pc, file, line, callerOk := runtime.Caller(1)

	err.code = 501
	err.message = "Unknown error"
	err.description = ""
	err.devCode = 1000
	err.devMessage = e.Error()
	err.devDescription = ""
	if callerOk {
		err.trace = file[strings.LastIndex(file, "/")+1:] +
		":" + runtime.FuncForPC(pc).Name() +
		":" + strconv.Itoa(line)
	} else {
		err.trace = ""
	}
	return err
}

func processError(e error) error {
	var err *MError
	var callerOk bool
	pc, file, line, callerOk := runtime.Caller(1)

	switch {
	case errors.As(e, &err):
		if err.trace == "" {
			err.trace = file[strings.LastIndex(file, "/")+1:] +
			":" + runtime.FuncForPC(pc).Name() +
			":" + strconv.Itoa(line)
		}
		if callerOk {
			err.trace = runtime.FuncForPC(pc).Name() + "/" + err.trace
		}
	default:
		err = err.New(e.Error())
	}
	return err
}
