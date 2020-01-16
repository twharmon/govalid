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

func (f32c *float32Constraint) error(val reflect.Value) error {
	f32 := val.Interface().(float32)
	if !f32c.req && f32 == 0 {
		return nil
	}
	if f32c.req && f32 == 0 {
		return fmt.Errorf("%w: %s is required", ErrInvalidStruct, f32c.field)
	}
	if f32c.isMaxSet && f32 > f32c.max {
		return fmt.Errorf("%w: %s can not be greater than %f", ErrInvalidStruct, f32c.field, f32c.max)
	}
	if f32c.isMinSet && f32 < f32c.min {
		return fmt.Errorf("%w: %s must be at least %f", ErrInvalidStruct, f32c.field, f32c.min)
	}
	return nil
}

func (f32c *float32Constraint) errors(val reflect.Value) []error {
	var vs []error
	f32 := val.Interface().(float32)
	if !f32c.req && f32 == 0 {
		return nil
	}
	if f32c.req && f32 == 0 {
		vs = append(vs, fmt.Errorf("%w: %s is required", ErrInvalidStruct, f32c.field))
	}
	if f32c.isMaxSet && f32 > f32c.max {
		vs = append(vs, fmt.Errorf("%w: %s can not be greater than %f", ErrInvalidStruct, f32c.field, f32c.max))
	}
	if f32c.isMinSet && f32 < f32c.min {
		vs = append(vs, fmt.Errorf("%w: %s must be at least %f", ErrInvalidStruct, f32c.field, f32c.min))
	}
	return vs
}
