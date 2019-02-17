package govalid

import "reflect"

// MustPrepare .
func MustPrepare(models ...interface{}) error {
	store = make(constraintMap)
	for _, model := range models {
		if err := mustPrepare(model); err != nil {
			return err
		}
	}
	return nil
}

func mustPrepare(model interface{}) error {
	typ := reflect.TypeOf(model).Elem()
	name := typ.Name()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			if err := makeStringConstraint(name, field); err != nil {
				return err
			}
		case reflect.Int:
			if err := makeIntConstraint(name, field); err != nil {
				return err
			}
		case reflect.Int64:
			if err := makeInt64Constraint(name, field); err != nil {
				return err
			}
		default:
			makeNilConstraint(name)
		}
	}
	return nil
}
