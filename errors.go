package errors

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime"
	"net/http"
)

const ERROR_UNKNOWN = "UnknownError"
const ERROR_UNAUTHORIZED = "Unauthorized"
const ERROR_NOT_FOUND = "NotFound"
const ERROR_INTERNAL = "ServerInternalError"
const ERROR_INVALID_PARAMS = "InvalidParams"
const ERROR_INTERNAL_EXCEPTION = "InternalException"
const ERROR_ALREADY_EXISTS = "AlreadyExists"

type ParamError struct {
	Field   string `json:"field,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Error struct {
	Status  int           `json:"status,omitempty"`
	Code    string        `json:"code,omitempty"`
	Message string        `json:"message,omitempty"`
	Errors  []*ParamError `json:"errors,omitempty"`
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
	return &Error{Status: http.StatusInternalServerError, Code: ERROR_UNKNOWN, Message: message}
}

func Wrap(err interface{}) *Error {
	if err == nil {
		panic("err is nil")
	}

	switch err.(type) {
	case *Error:
		return err.(*Error)
	case error:
		return &Error{
			Status:  http.StatusInternalServerError,
			Code:    ERROR_UNKNOWN,
			Message: err.(error).Error(),
		}
	default:
		return &Error{
			Status:  http.StatusInternalServerError,
			Code:    ERROR_UNKNOWN,
			Message: fmt.Sprint(err),
		}
	}
}

func InvalidParams(params ...*ParamError) *Error {
	return &Error{Status: http.StatusBadRequest, Code: ERROR_INVALID_PARAMS, Errors: params}
}

func InvalidParam(field string, message string) *Error {
	return &Error{
		Status: http.StatusBadRequest,
		Code:   ERROR_INVALID_PARAMS,
		Errors: []*ParamError{{Field: field, Code: ERROR_INVALID_PARAMS, Message: message}}}
}

func BadRequest(code string, message string) *Error {
	return &Error{Status: http.StatusBadRequest, Code: code, Message: message}
}

func Unauthorized(message string) *Error {
	return &Error{Status: http.StatusUnauthorized, Code: ERROR_UNAUTHORIZED, Message: message}
}

func NotFound(message string) *Error {
	return &Error{Status: http.StatusNotFound, Code: ERROR_NOT_FOUND, Message: message}
}

func AlreadyExists(message string) *Error {
	return &Error{Status: http.StatusBadRequest, Code: ERROR_ALREADY_EXISTS, Message: message}
}

func InternalServerError(message string) *Error {
	return &Error{Status: http.StatusInternalServerError, Code: ERROR_INTERNAL, Message: message}
}
