package govalid

import "reflect"

type nilConstraint struct{}

func (nc *nilConstraint) validate(val reflect.Value) []string {
	return nil
}

func makeNilConstraint(name string) {
	modelStore.add(name, new(nilConstraint))
}
