package govalid

import (
	"errors"
)

type ValidationError interface {
	// govalidError could be replaced in future with
	// helpful functions, like Field() string
	govalidError()
	Error() string
}

type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}

func (e *validationError) govalidError() {
	panic("do not call this")
}

func NewValidationError(msg string) ValidationError {
	return &validationError{msg: msg}
}

func wrap(prefix string, err error) error {
	verr, ok := err.(*validationError)
	if ok {
		return NewValidationError(prefix + ": " + verr.Error())
	}

	return errors.New(prefix + ": " + err.Error())
}

var _ error = (*validationError)(nil)
var _ ValidationError = (*validationError)(nil)
