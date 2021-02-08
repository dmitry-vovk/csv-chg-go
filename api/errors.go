package api

import (
	"errors"
	"strconv"
)

var (
	ErrBadRequest  = errors.New("bad request")
	ErrServerError = errors.New("internal server error")
)

// ErrUnexpectedStatusCode implements `error` with custom http status code
type ErrUnexpectedStatusCode struct {
	code int
}

func (e ErrUnexpectedStatusCode) Error() string {
	return "unexpected status code " + strconv.Itoa(e.code)
}

// ErrInvalidContentType implements `error` with custom http status code
type ErrInvalidContentType struct {
	contentType string
}

func (e ErrInvalidContentType) Error() string {
	return "invalid content type '" + e.contentType + "'"
}
