package govalid

import (
	"fmt"
	"reflect"
)

type float32Constraint struct {
	field    string
	req      bool
	isMinSet bool
	min      float32
	isMaxSet bool
	max      float32
}

func (f32c *float32Constraint) violation(val reflect.Value) string {
	f32 := val.Interface().(float32)
	if !f32c.req && f32 == 0 {
		return ""
	}
	if f32c.req && f32 == 0 {
		return fmt.Sprintf("%s is required", f32c.field)
	}
	if f32c.isMaxSet && f32 > f32c.max {
		return fmt.Sprintf("%s can not be greater than %f", f32c.field, f32c.max)
	}
	if f32c.isMinSet && f32 < f32c.min {
		return fmt.Sprintf("%s must be at least %f", f32c.field, f32c.min)
	}
	return ""
}

func (f32c *float32Constraint) violations(val reflect.Value) []string {
	var vs []string
	f32 := val.Interface().(float32)
	if !f32c.req && f32 == 0 {
		return nil
	}
	if f32c.req && f32 == 0 {
		vs = append(vs, fmt.Sprintf("%s is required", f32c.field))
	}
	if f32c.isMaxSet && f32 > f32c.max {
		vs = append(vs, fmt.Sprintf("%s can not be greater than %f", f32c.field, f32c.max))
	}
	if f32c.isMinSet && f32 < f32c.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %f", f32c.field, f32c.min))
	}
	return vs
}
