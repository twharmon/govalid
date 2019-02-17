package govalid

import (
	"fmt"
	"reflect"
)

// MustValidate .
func MustValidate(s interface{}) error {
	val := reflect.ValueOf(s).Elem()
	name := reflect.TypeOf(s).Elem().Name()
	m := store[name]
	if m == nil {
		return fmt.Errorf("%s was not prepared", name)
	}
	for i, c := range store[name] {
		if err := c.validate(val.Field(i)); err != nil {
			return err
		}
	}
	return nil
}
