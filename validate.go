package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

// Validate checks the struct s against all constraints and custom
// validation functions, if any. It returns a violation as as string.
// And empty string("") means there was no violation.
//
// An error is returned if the struct wasn't registered, if s is not
// a pointer to a struct, or if your custom validation functions
// return an error.
func Validate(s interface{}) (string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return "", errors.New("s must be a pointer to a struct")
	}
	e := t.Elem()
	if e.Kind() != reflect.Struct {
		return "", errors.New("s must be a pointer to a struct")
	}
	m := modelStore[e.Name()]
	if m == nil {
		return "", fmt.Errorf("%s was not registered", e.Name())
	}
	return m.validate(s)
}
