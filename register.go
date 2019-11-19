package govalid

import (
	"reflect"
	"strings"
)

// Register is required for all structs that you wish
// to validate. It is intended to be ran at load time
// and caches information about the structs to reduce
// run time allocations.
//
// NOTE: This is not thread safe. You must
// register structs before validating.
func Register(structs ...interface{}) {
	for _, s := range structs {
		register(s)
	}
}

func register(s interface{}) {
	typ := reflect.TypeOf(s)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		panic("only structs can be registered")
	}
	name := typ.Name()
	m := new(model)
	m.name = name
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		firstLetter := string(field.Name[0])
		if firstLetter != strings.ToUpper(firstLetter) {
			m.registerNilConstraint()
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			m.registerStringConstraint(field)
		case reflect.Int:
			m.registerIntConstraint(field)
		case reflect.Int64:
			m.registerInt64Constraint(field)
		case reflect.Float32:
			m.registerFloat32Constraint(field)
		case reflect.Float64:
			m.registerFloat64Constraint(field)
		case reflect.Struct:
			if _, ok := field.Type.FieldByName("String"); ok {
				m.registerStringConstraint(field)
			} else if _, ok := field.Type.FieldByName("Int64"); ok {
				m.registerInt64Constraint(field)
			} else if _, ok := field.Type.FieldByName("Float64"); ok {
				m.registerFloat64Constraint(field)
			} else if _, ok := field.Type.FieldByName("Time"); ok {
				m.registerTimeConstraint(field)
			} else if field.Type.String() == "time.Time" {
				m.registerTimeConstraint(field)
			} else {
				m.registerNilConstraint()
			}
		default:
			m.registerNilConstraint()
		}
	}
	m.addToRegistry(name)
}
