package govalid

import (
	"errors"
	"reflect"
)

// ErrNotStruct is encountered when an attempt is made to validate
// a type that is not a struct is made.
var ErrNotStruct = errors.New("only structs can be validated")

// ErrNotRegistered is encountered when an attempt is made to
// validate a type that has not yet been registered.
var ErrNotRegistered = errors.New("only structs can be validated")

// Violation checks the struct s against all constraints and custom
// validation functions, if any. It returns an error if the struct
// fails validation. If the type being validated is not a struct,
// ErrNotStruct will be returned. If the type being validated has not
// yet been registered, ErrNotRegistered is returned.
func Violation(s interface{}) error {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	m := modelStore[t.Name()]
	if m == nil {
		return ErrNotRegistered
	}
	return m.violation(s)
}

// Violations checks the struct s against all constraints and custom
// validation functions, if any. It returns a slice of errors if the
// struct fails validation. If the type being validated is not a
// struct, ErrNotStruct alone will be returned. If the type being
// validated has not yet been registered, ErrNotRegistered alone is
// returned.
func Violations(s interface{}) []error {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return []error{ErrNotStruct}
	}
	m := modelStore[t.Name()]
	if m == nil {
		return []error{ErrNotRegistered}
	}
	return m.violations(s)
}
