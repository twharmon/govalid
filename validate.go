package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

// Validate checks the struct s against all constraints and custom
// validation functions, if any. It returns violations as []string.
// An error is returned if the struct wasn't registered or if s
// is not a pointer to a struct.
func Validate(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New("s must be a pointer to a struct")
	}
	e := t.Elem()
	if e.Kind() != reflect.Struct {
		return nil, errors.New("s must be a pointer to a struct")
	}
	m := modelStore[e.Name()]
	if m == nil {
		return nil, fmt.Errorf("%s was not registered", e.Name())
	}
	return m.validate(s)
}
