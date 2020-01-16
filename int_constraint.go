package govalid

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type intConstraint struct {
	field    string
	req      bool
	isMinSet bool
	min      int
	isMaxSet bool
	max      int
	in       []int
}

func (ic *intConstraint) error(val reflect.Value) error {
	i := val.Interface().(int)
	if !ic.req && i == 0 {
		return nil
	}
	if ic.req && i == 0 {
		return fmt.Errorf("%w: %s is required", ErrInvalidStruct, ic.field)
	}
	if ic.isMaxSet && i > ic.max {
		return fmt.Errorf("%w: %s can not be greater than %d", ErrInvalidStruct, ic.field, ic.max)
	}
	if ic.isMinSet && i < ic.min {
		return fmt.Errorf("%w: %s must be at least %d", ErrInvalidStruct, ic.field, ic.min)
	}
	if len(ic.in) > 0 {
		for _, opt := range ic.in {
			if i == opt {
				return nil
			}
		}
	} else {
		return nil
	}
	iStrSlice := []string{}
	for _, a := range ic.in {
		iStrSlice = append(iStrSlice, strconv.Itoa(a))
	}
	return fmt.Errorf("%w: %s must be in [%s]", ErrInvalidStruct, ic.field, strings.Join(iStrSlice, ", "))
}

func (ic *intConstraint) errors(val reflect.Value) []error {
	var vs []error
	i := val.Interface().(int)
	if !ic.req && i == 0 {
		return nil
	}
	if ic.req && i == 0 {
		vs = append(vs, fmt.Errorf("%w: %s is required", ErrInvalidStruct, ic.field))
	}
	if ic.isMaxSet && i > ic.max {
		vs = append(vs, fmt.Errorf("%w: %s can not be greater than %d", ErrInvalidStruct, ic.field, ic.max))
	}
	if ic.isMinSet && i < ic.min {
		vs = append(vs, fmt.Errorf("%w: %s must be at least %d", ErrInvalidStruct, ic.field, ic.min))
	}
	if len(ic.in) > 0 {
		for _, opt := range ic.in {
			if i == opt {
				return vs
			}
		}
	} else {
		return vs
	}
	iStrSlice := []string{}
	for _, a := range ic.in {
		iStrSlice = append(iStrSlice, strconv.Itoa(a))
	}
	return append(vs, fmt.Errorf("%w: %s must be in [%s]", ErrInvalidStruct, ic.field, strings.Join(iStrSlice, ", ")))
}
