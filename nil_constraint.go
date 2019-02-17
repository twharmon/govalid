package govalid

import "reflect"

type nilConstraint struct{}

func (nc *nilConstraint) validate(val reflect.Value) error {
	return nil
}

func makeNilConstraint(name string) {
	constraintStore.Add(name, new(nilConstraint))
}
