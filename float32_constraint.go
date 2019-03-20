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

func (f32c *float32Constraint) validate(val reflect.Value) string {
	empty := true
	f32, ok := val.Interface().(float32)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		f32 = val.FieldByName("Float32").Interface().(float32)
		empty = false
	} else {
		empty = f32 == 0
	}
	if !f32c.req && empty {
		return ""
	}
	if f32c.req && empty {
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
