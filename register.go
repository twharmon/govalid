package govalid

import "reflect"

// Register is required for all structs that you wish
// to validate. It is intended to be ran at load time
// and cashes information about the structs to reduce
// run time allocations.
//
// NOTE: This is not thread safe. You must
// register structs before validating.
func Register(structs ...interface{}) error {
	constraintStore = make(constraintMap)
	for _, s := range structs {
		if err := register(s); err != nil {
			return err
		}
	}
	return nil
}

func register(s interface{}) error {
	typ := reflect.TypeOf(s).Elem()
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
