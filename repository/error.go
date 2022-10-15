package repository

import (
	"errors"
	"net/http"
)

type MyError struct {
	msg  string
	code int
}

func NewError(msg string) *MyError {
	return &MyError{
		msg:  msg,
		code: http.StatusBadRequest,
	}
}

func (e *MyError) Error() string {
	return e.msg
}

// Code returns http status code
func (e *MyError) Code() int {
	return e.code
}

var ErrInternalServerError error = errors.New("Internal Server Error")
