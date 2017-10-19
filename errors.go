package errors

import (
	"fmt"
	"net/http"
)

const ERROR_UNKNOWN = "UnknownError"
const ERROR_UNAUTHORIZED = "Unauthorized"
const ERROR_NOT_FOUND = "NotFound"
const ERROR_INTERNAL = "ServerInternalError"
const ERROR_INVALID_PARAMS = "InvalidParams"
const ERROR_INTERNAL_EXCEPTION = "InternalException"

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

func (e *Error) Error() string {
	return fmt.Sprint(e.Status, e.Code, e.Message, e.Errors)
}

func InvalidParams(params ...*ParamError) *Error {
	return &Error{Status: http.StatusBadRequest, Code: ERROR_INVALID_PARAMS, Errors: params}
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

func InternalServerError(message string) *Error {
	return &Error{Status: http.StatusInternalServerError, Code: ERROR_INTERNAL, Message: message}
}