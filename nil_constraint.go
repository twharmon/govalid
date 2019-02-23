package govalid

import "reflect"

type nilConstraint struct{}

func (nc *nilConstraint) validate(val reflect.Value) string {
	return ""
}
