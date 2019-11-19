package govalid

import (
	"log"
	"reflect"
)

// Violation checks the struct s against all constraints and custom
// validation functions, if any. It returns an error if the struct
// fails validation.
func Violation(s interface{}) error {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("s must be a pointer to a struct")
		return nil
	}
	m := modelStore[t.Name()]
	if m == nil {
		log.Printf("%s was not registered\n", t.Name())
		return nil
	}
	return m.violation(s)
}

// Violations checks the struct s against all constraints and custom
// validation functions, if any. It returns a slice of errors if the
// struct fails validation.
func Violations(s interface{}) []error {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("s must be a pointer to a struct")
		return nil
	}
	m := modelStore[t.Name()]
	if m == nil {
		log.Printf("%s was not registered\n", t.Name())
		return nil
	}
	return m.violations(s)
}
