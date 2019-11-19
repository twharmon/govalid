package govalid

import (
	"reflect"
)

type constraint interface {
	violation(reflect.Value) error
	violations(reflect.Value) []error
}

// AddCustom adds custom validation functions to struct s.
//
// NOTE: This is not thread safe. You must
// add cusrom validation functions before validating.
func AddCustom(s interface{}, f ...func(interface{}) error) {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	n := t.Name()
	m := modelStore[n]
	if m == nil {
		panic("struct s must be registered before adding a custom validator")
	}
	m.custom = append(m.custom, f...)
}
