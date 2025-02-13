package govalid

import (
	"fmt"
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
		return NewValidationError(fmt.Sprintf("%s: %s", prefix, verr))
	}
	return fmt.Errorf("%s: %w", prefix, err)
}

var _ error = (*validationError)(nil)
var _ ValidationError = (*validationError)(nil)
