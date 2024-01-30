// Package merror provides utilities for handling and rendering errors in a Gin-based application.

package merror

import (
	"fmt"
	"strconv"
	"errors"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc/metadata"
)

// MultipassError represents a structured error format with code, message, and description.
type MultipassError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// RenderError handles and renders errors in the Gin context.
// It sets the HTTP status, sends an HTTP error message, and logs developer error details.
func RenderError(c *gin.Context, err error) {
	var mError *MError
	HTTPError := MultipassError{}
	colorYellow := "\u001b[33m"
	colorNone := "\033[0m"

	// Check if the error is of type MError.
	switch {
	case errors.As(err, &mError):
		// If it is, set HTTP status and render HTTP error.
		HTTPError.Code = mError.Code
		HTTPError.Message = mError.Message
		HTTPError.Description = mError.Description
		c.AbortWithStatus(422)
		errorJSON, _ := json.Marshal(HTTPError)
		c.Data(mError.Code, "application/json", errorJSON)

		// Log error for developer.
		fmt.Println(string(colorYellow), "->Developer Error:")
		fmt.Println("Raised at: ", time.Now().String())
		fmt.Println("Code: ", mError.DevCode)
		fmt.Println("Message: ", mError.DevMessage)
		fmt.Println("Description: ", mError.DevDescription)
		fmt.Println("Trace: ", mError.Trace, string(colorNone))
	default:
		// If not, set a generic 501 status and render an unknown error.
		c.AbortWithStatus(501)
		HTTPError.Code = 501
		HTTPError.Message = "Unknown Error"
		HTTPError.Description = ""
	}
}

// ToMError converts gRPC metadata trailer to a MultipassError.
func ToMError(trailer metadata.MD) error {
	mError := &MError{}
	if code, ok := trailer["code"]; ok {
		c, err := strconv.Atoi(code[0])
		if err != nil {
			// If code conversion fails, set a default code.
			mError.Code = 422
		} else {
			mError.Code = c
		}
	}
	if message, ok := trailer["message"]; ok {
		mError.Message = message[0]
	}
	if description, ok := trailer["description"]; ok {
		mError.Description = description[0]
	}
	if trace, ok := trailer["trace"]; ok {
		mError.Trace = trace[0]
	}
	if DevMessage, ok := trailer["dev_message"]; ok {
		mError.DevMessage = DevMessage[0]
	}
	if DevDescription, ok := trailer["dev_description"]; ok {
		mError.DevDescription = DevDescription[0]
	}
	if code, ok := trailer["code"]; ok {
		c, err := strconv.Atoi(code[0])
		if err != nil {
			// If code conversion fails, set a default code.
			mError.Code = 422
		} else {
			mError.Code = c
		}
	}
	if devCode, ok := trailer["dev_code"]; ok {
		c, err := strconv.Atoi(devCode[0])
		if err != nil {
			// If dev_code conversion fails, set a default dev_code.
			mError.DevCode = 501
		} else {
			mError.DevCode = c
		}
	}

	return mError
}