package errors

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

type Error struct {
	Message string `json:"message" description:"error message"`
}

var None = Error{Message: "success"}

func (e Error) Error() string {
	return e.Message
}

func Wrap(err error) error {
	return Error{Message: err.Error()}
}

func New(format string, args ...interface{}) error {
	return Error{Message: fmt.Sprintf(format, args...)}
}

func GetServiceErrorCode(err error) int {
	if svcErr, ok := err.(restful.ServiceError); ok {
		return svcErr.Code
	} else {
		return http.StatusInternalServerError
	}
}
