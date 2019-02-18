package govalid

import (
	"fmt"
	"reflect"
)

// Validate .
func Validate(s interface{}) error {
	val := reflect.ValueOf(s).Elem()
	name := reflect.TypeOf(s).Elem().Name()
	constraints := constraintStore[name]
	if constraints == nil {
		return fmt.Errorf("%s was not prepared", name)
	}
	for i, c := range constraints {
		if err := c.validate(val.Field(i)); err != nil {
			return err
		}
	}
	return nil
}
