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
	if t.Kind() != reflect.Ptr {
		log.Println("s must be a pointer to a struct")
		return nil
	}
	e := t.Elem()
	if e.Kind() != reflect.Struct {
		log.Println("s must be a pointer to a struct")
		return nil
	}
	m := modelStore[e.Name()]
	if m == nil {
		log.Printf("%s was not registered\n", e.Name())
		return nil
	}
	return m.violation(s)
}

// Violations checks the struct s against all constraints and custom
// validation functions, if any. It returns an error if the struct
// fails validation.
func Violations(s interface{}) []error {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		log.Println("s must be a pointer to a struct")
		return nil
	}
	e := t.Elem()
	if e.Kind() != reflect.Struct {
		log.Println("s must be a pointer to a struct")
		return nil
	}
	m := modelStore[e.Name()]
	if m == nil {
		log.Printf("%s was not registered\n", e.Name())
		return nil
	}
	return m.violations(s)
}
