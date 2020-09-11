package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// AbortWithError stops the chain, writes the status code and the given error
func AbortWithError(ctx *gin.Context, code int, err error, details string) {
	ctx.AbortWithStatusJSON(code, &Error{
		Status:  http.StatusText(code),
		Code:    code,
		Message: err.Error(),
		Details: details,
	})
}

func NewNotFoundError(message, details string) *Error {
	return &Error{
		Status:  http.StatusText(http.StatusNotFound),
		Code:    http.StatusNotFound,
		Message: message,
		Details: details,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s %s: %s", e.Code, e.Status, e.Message, e.Details)
}

func IsNotFound(err error) bool {
	switch t := err.(type) {
	case Error:
		return t.Code == http.StatusNotFound
	case *Error:
		return t.Code == http.StatusNotFound
	}
	return false
}
