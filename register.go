package govalid

import (
	"reflect"
)

// Register is required for all structs that you wish
// to validate. It is intended to be ran at load time
// and cashes information about the structs to reduce
// run time allocations.
//
// NOTE: This is not thread safe. You must
// register structs before validating.
func Register(structs ...interface{}) {
	modelStore = make(modelMap)
	for _, s := range structs {
		register(s)
	}
}

func register(s interface{}) {
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		panic("pointers can not be registered")
	}
	if typ.Kind() != reflect.Struct {
		panic("only structs can be registered")
	}
	name := typ.Name()
	modelStore[name] = new(model)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			makeStringConstraint(name, field)
		case reflect.Int:
			makeIntConstraint(name, field)
		case reflect.Int64:
			makeInt64Constraint(name, field)
		default:
			makeNilConstraint(name)
		}
	}
}
