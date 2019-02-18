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
	val := reflect.ValueOf(s).Elem()
	name := t.Elem().Name()
	model := modelStore[name]
	if model == nil {
		return nil, fmt.Errorf("%s was not registered", name)
	}
	var vs []string
	for i, c := range model.constraints {
		vs = append(vs, c.validate(val.Field(i))...)
	}
	for _, v := range model.custom {
		vs = append(vs, v(s)...)
	}
	return vs, nil
}
