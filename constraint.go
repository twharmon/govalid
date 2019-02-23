package govalid

import (
	"reflect"
)

type constraint interface {
	validate(reflect.Value) string
}

// AddCustom adds custom validation functions to struct s.
//
// NOTE: This is not thread safe. You must
// add cusrom validation functions before validating.
func AddCustom(s interface{}, f ...func(interface{}) (string, error)) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		panic("s can not be a pointer")
	}
	n := t.Name()
	m := modelStore[n]
	if m == nil {
		panic("struct s must be registered before adding a custom validator")
	}
	m.custom = append(m.custom, f...)
}
