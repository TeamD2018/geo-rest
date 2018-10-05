package models

import "fmt"

type Error struct {
	Parameter interface{} `json:"-"`
	Message   string      `json:"msg"`
	Code      int         `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.Code, fmt.Sprintf(e.Message, e.Parameter))
}

func (e *Error) SetParameter(p interface{}) *Error {
	e.Parameter = p
	return e
}

// List of Error
var (
	ErrCourierAlreadyExists              = Error{Message: "Courier with id %s already exists", Code: 10}
	ErrServerError                       = Error{Message: "Sorry, server error", Code: 20}
	ErrCourierNotFound                   = Error{Message: "Courier with id %s not found", Code: 60}
	ErrOneOfParametersNotFound           = Error{Message: "One of parameters not found", Code: 30}
	ErrOneOfParameterHaveIncorrectFormat = Error{Message: "One of parameter (%s) have incorrect format", Code: 40}
	ErrEntityNotFound                    = Error{Message: "Entity with such id %s not found", Code: 50}
	ErrUnmarshalJSON                     = Error{Message: "Error with unmarshal JSON: %s"}
)
