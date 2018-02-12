package errors

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	openapiErrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"net/http"
)

const ErrUnknown = "Unknown"
const ErrUnauthorized = "Unauthorized"
const ErrNotFound = "NotFound"
const ErrInvalidParams = "InvalidParams"
const ErrAlreadyExists = "AlreadyExists"

type Error struct {
	Status  int    `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(e.Status)

	enc := json.NewEncoder(rw)
	err := enc.Encode(e)
	if err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

func (e *Error) Error() string {
	s, _ := json.Marshal(e)
	return string(s)
}

func Unknown(message string) *Error {
	return &Error{Status: http.StatusInternalServerError, Code: ErrUnknown, Message: message}
}

func InvalidParam(message string) *Error {
	return &Error{Status: http.StatusBadRequest, Code: ErrInvalidParams, Message: message}
}

func BadRequest(code string, message string) *Error {
	return &Error{Status: http.StatusBadRequest, Code: code, Message: message}
}

func Unauthorized(message string) *Error {
	return &Error{Status: http.StatusUnauthorized, Code: ErrUnauthorized, Message: message}
}

func NotFound(message string) *Error {
	return &Error{Status: http.StatusNotFound, Code: ErrNotFound, Message: message}
}

func AlreadyExists(message string) *Error {
	return &Error{Status: http.StatusBadRequest, Code: ErrAlreadyExists, Message: message}
}

func Wrap(err interface{}) *Error {
	if err == nil {
		panic("err is nil")
	}

	switch err.(type) {
	case *Error:
		return err.(*Error)
	case *openapiErrors.MethodNotAllowedError:
		return &Error{Status: http.StatusMethodNotAllowed, Code: "MethodNotAllowed", Message: "MethodNotAllowed"}
	case *jwt.ValidationError:
		return Unauthorized(err.(*jwt.ValidationError).Error())
	case openapiErrors.Error:
		if err == nil {
			return Unknown("e==nil")
		}

		e := err.(openapiErrors.Error)
		return &Error{Status: int(e.Code()), Code: ErrUnknown, Message: e.Error()}
	case error:
		return Unknown(err.(error).Error())
	default:
		return Unknown(fmt.Sprint(err))
	}
}
