package govalid

import "reflect"

type nilConstraint struct{}

func (nc *nilConstraint) violation(val reflect.Value) string {
	return ""
}

func (nc *nilConstraint) violations(val reflect.Value) []string {
	return nil
}
