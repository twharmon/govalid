package govalid

import (
	"errors"
)

// ErrNotStruct is encountered when an attempt is made to validate
// a type that is not a struct is made.
var ErrNotStruct = errors.New("only structs can be validated")

// Deprecated: ErrNotRegistered is deprecated.
var ErrNotRegistered = errors.New("structs must be registered before validating")

// New .
func New() *Validator {
	v := new(Validator)
	v.store = make(map[string]*model)
	v.stringRules = make(map[string]func(string, string) string)
	v.int64Rules = make(map[string]func(string, int64) string)
	v.float64Rules = make(map[string]func(string, float64) string)
	return v
}
