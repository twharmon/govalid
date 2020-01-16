package govalid

import "reflect"

type nilConstraint struct{}

func (nc *nilConstraint) error(val reflect.Value) error {
	return nil
}

func (nc *nilConstraint) errors(val reflect.Value) []error {
	return nil
}
