package govalid

import "reflect"

type constraintMap map[string][]constraint

var constraintStore constraintMap

func (cs constraintMap) Add(model string, c constraint) {
	cs[model] = append(cs[model], c)
}

type constraint interface {
	validate(reflect.Value) error
}
