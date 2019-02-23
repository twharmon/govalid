package govalid

import (
	"fmt"
	"reflect"
)

type float64Constraint struct {
	field string
	req   bool
	min   float64
	max   float64
}

func (f64c *float64Constraint) validate(val reflect.Value) string {
	f64 := val.Interface().(float64)
	if !f64c.req && f64 == 0 {
		return ""
	}
	if f64c.req && f64 == 0 {
		return fmt.Sprintf("%s is required", f64c.field)
	}
	if f64c.max > 0 && f64 > f64c.max {
		return fmt.Sprintf("%s can not be greater than %f", f64c.field, f64c.max)
	}
	if f64c.min > 0 && f64 < f64c.min {
		return fmt.Sprintf("%s must be at least %f", f64c.field, f64c.min)
	}
	return ""
}
