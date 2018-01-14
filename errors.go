package errors

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime"
	"net/http"
)

const ErrUnknown = "Unknown"
const ErrUnauthorized = "Unauthorized"
const ErrNotFound = "NotFound"
const ErrInvalidParams = "InvalidParams"
const ErrAlreadyExists = "AlreadyExists"

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
			Code:    ErrUnknown,
			Message: err.(error).Error(),
		}
	default:
		return &Error{
			Status:  http.StatusInternalServerError,
			Code:    ErrUnknown,
			Message: fmt.Sprint(err),
		}
	}
}

func Unknown(message string) *Error {
	return &Error{Status: http.StatusInternalServerError, Code: ErrUnknown, Message: message}
}

func InvalidParams(params ...*ParamError) *Error {
	return &Error{Status: http.StatusBadRequest, Code: ErrInvalidParams, Errors: params}
}

func InvalidParam(field string, message string) *Error {
	return &Error{
		Status: http.StatusBadRequest,
		Code:   ErrInvalidParams,
		Errors: []*ParamError{{Field: field, Code: ErrInvalidParams, Message: message}}}
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
