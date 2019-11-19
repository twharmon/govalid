package govalid

import (
	"errors"
	"reflect"
)

// ErrNotPtrToStruct is encountered when an attempt is made to validate
// a type that is not a struct is made.
var ErrNotPtrToStruct = errors.New("only pointers to structs can be validated")

// ErrNotRegistered is encountered when an attempt is made to
// validate a type that has not yet been registered.
var ErrNotRegistered = errors.New("structs must be registered before validating")

// Register is required for all structs that you wish
// to validate. It is intended to be ran at load time
// and caches information about the structs to reduce
// run time allocations.
//
// NOTE: This is not thread safe. You must
// register structs before validating.
func Register(structs ...interface{}) {
	for _, s := range structs {
		register(s)
	}
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

// Violation checks the struct s against all constraints and custom
// validation functions, if any. It returns an error if the struct
// fails validation. If the type being validated is not a struct,
// ErrNotPtrToStruct will be returned. If the type being validated
// has not yet been registered, ErrNotRegistered is returned.
func Violation(s interface{}) error {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return ErrNotPtrToStruct
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return ErrNotPtrToStruct
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
// struct, ErrNotPtrToStruct alone will be returned. If the type
// being validated has not yet been registered, ErrNotRegistered
// alone is returned.
func Violations(s interface{}) []error {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return []error{ErrNotPtrToStruct}
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return []error{ErrNotPtrToStruct}
	}
	m := modelStore[t.Name()]
	if m == nil {
		return []error{ErrNotRegistered}
	}
	return m.violations(s)
}
