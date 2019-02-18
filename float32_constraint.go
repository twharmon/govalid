package govalid

import (
	"fmt"
	"reflect"
)

type float32Constraint struct {
	field string
	req   bool
	min   float32
	max   float32
}

func (f32c *float32Constraint) validate(val reflect.Value) []string {
	f32 := val.Interface().(float32)
	var vs []string
	if !f32c.req && f32 == 0 {
		return vs
	}
	if f32c.req && f32 == 0 {
		vs = append(vs, fmt.Sprintf("%s is required", f32c.field))
		return vs
	}
	if f32c.max > 0 && f32 > f32c.max {
		vs = append(vs, fmt.Sprintf("%s can not be greater than %f", f32c.field, f32c.max))
	}
	if f32c.min > 0 && f32 < f32c.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %f", f32c.field, f32c.min))
	}
	return vs
}
