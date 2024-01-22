package merror

import (
	"fmt"
	"strconv"
	"errors"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc/metadata"
	me "github.com/rgrente/mtp-golibs/merror"
)

type MultipassError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// Raise HTTP Fail with specified status
// Send HTTP error message
// Log Developer error message (with more detailed error)
func renderError(c *gin.Context, err error) {
	var mError *me.MError
	HTTPError := MultipassError{}
	colorYellow := "\u001b[33m"
	colorNone := "\033[0m"

	switch {
	case errors.As(err, &mError):
		// Raise HTTP status and render HTTP error
		HTTPError.Code = mError.code
		HTTPError.Message = mError.message
		HTTPError.Description = mError.description
		c.AbortWithStatus(mError.code)
		errorJSON, _ := json.Marshal(HTTPError)
		c.Data(mError.code, "application/json", errorJSON)

		// Log error for developer
		fmt.Println(string(colorYellow), "Developer Error:")
		fmt.Println("Code: ", mError.devCode)
		fmt.Println("Message: ", mError.devMessage)
		fmt.Println("Description: ", mError.devDescription)
		fmt.Println("Trace: ", mError.trace, string(colorNone))
	default:
		c.AbortWithStatus(501)
		HTTPError.Code = 501
		HTTPError.Message = "Unknown Error"
		HTTPError.Description = ""
	}
}

func toMError(trailer metadata.MD) error {
	code, err := strconv.Atoi(trailer["code"][0])
	mError := &me.MError{}
	if err != nil {
		mError.code = 422
	} else {
		mError.code = code
	}
	devCode, err := strconv.Atoi(trailer["dev_code"][0])
	if err != nil {
		mError.devCode = 501
	} else {
		mError.devCode = devCode
	}

	mError.message = trailer["message"][0]
	mError.description = trailer["description"][0]
	mError.trace = trailer["trace"][0]
	mError.devMessage = trailer["dev_message"][0]
	mError.devDescription = trailer["dev_description"][0]

	return mError
}
