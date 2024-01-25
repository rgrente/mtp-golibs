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

type MultipassError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// Raise HTTP Fail with specified status
// Send HTTP error message
// Log Developer error message (with more detailed error)
func RenderError(c *gin.Context, err error) {
	var mError *MError
	HTTPError := MultipassError{}
	colorYellow := "\u001b[33m"
	colorNone := "\033[0m"

	switch {
	case errors.As(err, &mError):
		// Raise HTTP status and render HTTP error
		HTTPError.Code = 422
		HTTPError.Message = mError.Message
		HTTPError.Description = mError.Description
		c.AbortWithStatus(HTTPError.Code)
		errorJSON, _ := json.Marshal(HTTPError)
		c.Data(mError.Code, "application/json", errorJSON)

		// Log error for developer
		fmt.Println(string(colorYellow), "Developer Error:")
		fmt.Println("Raised at: ", time.Now().String())
		fmt.Println("Code: ", mError.DevCode)
		fmt.Println("Message: ", mError.DevMessage)
		fmt.Println("Description: ", mError.DevDescription)
		fmt.Println("Trace: ", mError.Trace, string(colorNone))
	default:
		c.AbortWithStatus(501)
		HTTPError.Code = 501
		HTTPError.Message = "Unknown Error"
		HTTPError.Description = ""
	}
}

func ToMError(trailer metadata.MD) error {
	code, err := strconv.Atoi(trailer["code"][0])
	mError := &MError{}
	if err != nil {
		mError.Code = 422
	} else {
		mError.Code = code
	}
	devCode, err := strconv.Atoi(trailer["dev_code"][0])
	if err != nil {
		mError.DevCode = 501
	} else {
		mError.DevCode = devCode
	}

	mError.Message = trailer["message"][0]
	mError.Description = trailer["description"][0]
	mError.Trace = trailer["trace"][0]
	mError.DevMessage = trailer["dev_message"][0]
	mError.DevDescription = trailer["dev_description"][0]

	return mError
}
