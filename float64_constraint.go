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

func (f64c *float64Constraint) violation(val reflect.Value) string {
	var empty bool
	f64, ok := val.Interface().(float64)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		f64 = val.FieldByName("Float64").Interface().(float64)
		empty = false
	} else {
		empty = f64 == 0
	}
	if !f64c.req && empty {
		return ""
	}
	if f64c.req && empty {
		return fmt.Sprintf("%s is required", f64c.field)
	}
	if f64c.isMaxSet && f64 > f64c.max {
		return fmt.Sprintf("%s can not be greater than %f", f64c.field, f64c.max)
	}
	if f64c.isMinSet && f64 < f64c.min {
		return fmt.Sprintf("%s must be at least %f", f64c.field, f64c.min)
	}
	return ""
}

func (f64c *float64Constraint) violations(val reflect.Value) []string {
	var vs []string
	var empty bool
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
		vs = append(vs, fmt.Sprintf("%s is required", f64c.field))
	}
	if f64c.isMaxSet && f64 > f64c.max {
		vs = append(vs, fmt.Sprintf("%s can not be greater than %f", f64c.field, f64c.max))
	}
	if f64c.isMinSet && f64 < f64c.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %f", f64c.field, f64c.min))
	}
	return vs
}
