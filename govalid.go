package govalid

import (
	"errors"
)

// ErrNotStruct is encountered when an attempt is made to validate
// a type that is not a struct is made.
var ErrNotStruct = errors.New("only structs can be validated")

// ErrNotRegistered is encountered when an attempt is made to
// validate a type that has not yet been registered.
var ErrNotRegistered = errors.New("structs must be registered before validating")

// New .
func New() *Validator {
	v := new(Validator)
	v.modelStore = make(map[string]*model)
	return v
}
