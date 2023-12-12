package apperr

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorCode int

const (
	ErrorCodeUnexpected ErrorCode = iota + 1
	ErrorCodeBadRequest
	ErrorCodeUnAuthorized
	ErrorCodeForbidden
	ErrorCodeNotFound
	ErrorCodeInternalServerError
)

type AppErr interface {
	Code() ErrorCode
	Message() string
	HTTPError() *echo.HTTPError
}

func (e *appErr) Wrap(err error) *appErr {
	return &appErr{
		code:    e.Code(),
		message: err.Error(),
	}
}

func (e *appErr) SetMessage(message string) *appErr {
	return &appErr{
		code:    e.Code(),
		message: message,
	}
}

type appErr struct {
	code    ErrorCode
	message string
}

func new(code ErrorCode, message string) *appErr {
	return &appErr{code, message}
}

func (e *appErr) Code() ErrorCode {
	return e.code
}

func (e *appErr) Message() string {
	return e.message
}

func (e *appErr) HTTPError() *echo.HTTPError {
	var code int
	switch e.Code() {
	case ErrorCodeBadRequest:
		code = http.StatusBadRequest
	case ErrorCodeUnAuthorized:
		code = http.StatusUnauthorized
	case ErrorCodeForbidden:
		code = http.StatusForbidden
	case ErrorCodeNotFound:
		code = http.StatusNotFound
	case ErrorCodeInternalServerError:
		code = http.StatusInternalServerError
	default:
		code = http.StatusInternalServerError
	}

	return &echo.HTTPError{
		Code:    code,
		Message: e.Message(),
	}
}

func NewUnexpectedError() *appErr {
	return new(ErrorCodeUnexpected, "unexpected error")
}

func NewBadRequestError() *appErr {
	return new(ErrorCodeBadRequest, "bad request")
}

func NewUnAuthorizedError() *appErr {
	return new(ErrorCodeUnAuthorized, "unauthorized")
}

func NewForbiddenError() *appErr {
	return new(ErrorCodeForbidden, "forbidden")
}

func NewNotFoundError() *appErr {
	return new(ErrorCodeNotFound, "not found")
}

func NewInternalServerError() *appErr {
	return new(ErrorCodeInternalServerError, "internal server error")
}
