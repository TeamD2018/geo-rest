package models

import (
	"fmt"
	"net/http"
)

type Error struct {
	Parameter interface{} `json:"-"`
	HttpCode  int         `json:"-"`
	Message   string      `json:"msg"`
	Code      int         `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.Code, fmt.Sprintf(e.Message, e.Parameter))
}

func (e *Error) HttpStatus() int {
	return e.HttpCode
}

func (e *Error) SetParameter(p interface{}) *Error {
	e.Parameter = p
	return e
}

// List of Error
var (
	ErrCourierAlreadyExists              = Error{Message: "Courier with id %s already exists", Code: 10, HttpCode: http.StatusConflict}
	ErrServerError                       = Error{Message: "Sorry, server error", Code: 20, HttpCode: http.StatusInternalServerError}
	ErrCourierNotFound                   = Error{Message: "Courier with id %s not found", Code: 60, HttpCode: http.StatusNotFound}
	ErrOneOfParametersNotFound           = Error{Message: "One of parameters not found", Code: 30, HttpCode: http.StatusBadRequest}
	ErrOneOfParameterHaveIncorrectFormat = Error{Message: "One of parameter (%s) have incorrect format", Code: 40, HttpCode: http.StatusBadRequest}
	ErrEntityNotFound                    = Error{Message: "Entity with such id %s not found", Code: 50, HttpCode: http.StatusNotFound}
	ErrUnmarshalJSON                     = Error{Message: "Error with unmarshal JSON: %s", Code: 70, HttpCode: http.StatusBadRequest}
)
