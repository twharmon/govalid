package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

// Validate checks the struct s against all constraints and custom
// validation functions, if any.
func Validate(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New("s must be a pointer")
	}
	name := t.Elem().Name()
	m := modelStore[name]
	if m == nil {
		return nil, fmt.Errorf("%s was not registered", name)
	}
	return m.validate(s), nil
}
