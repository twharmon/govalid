package govalid

import (
	"reflect"
)

type constraint interface {
	validate(reflect.Value) []string
}

type model struct {
	constraints []constraint
	custom      []func(interface{}) []string
}

type modelMap map[string]*model

var modelStore modelMap

func (mm modelMap) add(model string, c constraint) {
	mm[model].constraints = append(mm[model].constraints, c)
}

// AddCustom adds custom validation functions to struct s
func AddCustom(s interface{}, f func(interface{}) []string) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		panic("pointers can not be registered")
	}
	n := t.Name()
	m := modelStore[n]
	if m == nil {
		panic("struct s must be registered before adding a custom validator")
	}
	m.custom = append(m.custom, f)
}
