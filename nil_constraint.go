package govalid

import "reflect"

type nilConstraint struct{}

func (nc *nilConstraint) violation(val reflect.Value) error {
	return nil
}

func (nc *nilConstraint) violations(val reflect.Value) []error {
	return nil
}
