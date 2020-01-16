package govalid

import (
	"fmt"
	"reflect"
)

type float64Constraint struct {
	field    string
	req      bool
	isMinSet bool
	min      float64
	isMaxSet bool
	max      float64
}

func (f64c *float64Constraint) error(val reflect.Value) error {
	empty := true
	f64, ok := val.Interface().(float64)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		f64 = val.FieldByName("Float64").Interface().(float64)
		empty = false
	} else {
		empty = f64 == 0
	}
	if !f64c.req && empty {
		return nil
	}
	if f64c.req && empty {
		return fmt.Errorf("%w: %s is required", ErrInvalidStruct, f64c.field)
	}
	if f64c.isMaxSet && f64 > f64c.max {
		return fmt.Errorf("%w: %s can not be greater than %f", ErrInvalidStruct, f64c.field, f64c.max)
	}
	if f64c.isMinSet && f64 < f64c.min {
		return fmt.Errorf("%w: %s must be at least %f", ErrInvalidStruct, f64c.field, f64c.min)
	}
	return nil
}

func (f64c *float64Constraint) errors(val reflect.Value) []error {
	var vs []error
	empty := true
	f64, ok := val.Interface().(float64)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		f64 = val.FieldByName("Float64").Interface().(float64)
		empty = false
	} else {
		empty = f64 == 0
	}
	if !f64c.req && empty {
		return nil
	}
	if f64c.req && empty {
		vs = append(vs, fmt.Errorf("%w: %s is required", ErrInvalidStruct, f64c.field))
	}
	if f64c.isMaxSet && f64 > f64c.max {
		vs = append(vs, fmt.Errorf("%w: %s can not be greater than %f", ErrInvalidStruct, f64c.field, f64c.max))
	}
	if f64c.isMinSet && f64 < f64c.min {
		vs = append(vs, fmt.Errorf("%w: %s must be at least %f", ErrInvalidStruct, f64c.field, f64c.min))
	}
	return vs
}
