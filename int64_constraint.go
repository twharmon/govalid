package govalid

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type int64Constraint struct {
	field    string
	req      bool
	isMinSet bool
	min      int64
	isMaxSet bool
	max      int64
	in       []int64
}

func (i64c *int64Constraint) violation(val reflect.Value) error {
	empty := true
	i64, ok := val.Interface().(int64)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		i64 = val.FieldByName("Int64").Interface().(int64)
		empty = false
	} else {
		empty = i64 == 0
	}
	if !i64c.req && empty {
		return nil
	}
	if i64c.req && empty {
		return fmt.Errorf("%s is required", i64c.field)
	}
	if i64c.isMaxSet && i64 > i64c.max {
		return fmt.Errorf("%s can not be greater than %d", i64c.field, i64c.max)
	}
	if i64c.isMinSet && i64 < i64c.min {
		return fmt.Errorf("%s must be at least %d", i64c.field, i64c.min)
	}
	if len(i64c.in) > 0 {
		for _, opt := range i64c.in {
			if i64 == opt {
				return nil
			}
		}
	} else {
		return nil
	}
	iStrSlice := []string{}
	for _, a := range i64c.in {
		iStrSlice = append(iStrSlice, strconv.FormatInt(a, 10))
	}
	return fmt.Errorf("%s must be in [%s]", i64c.field, strings.Join(iStrSlice, ", "))
}

func (i64c *int64Constraint) violations(val reflect.Value) []error {
	var vs []error
	empty := true
	i64, ok := val.Interface().(int64)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		i64 = val.FieldByName("Int64").Interface().(int64)
		empty = false
	} else {
		empty = i64 == 0
	}
	if !i64c.req && empty {
		return nil
	}
	if i64c.req && empty {
		vs = append(vs, fmt.Errorf("%s is required", i64c.field))
	}
	if i64c.isMaxSet && i64 > i64c.max {
		vs = append(vs, fmt.Errorf("%s can not be greater than %d", i64c.field, i64c.max))
	}
	if i64c.isMinSet && i64 < i64c.min {
		vs = append(vs, fmt.Errorf("%s must be at least %d", i64c.field, i64c.min))
	}
	if len(i64c.in) > 0 {
		for _, opt := range i64c.in {
			if i64 == opt {
				return vs
			}
		}
	} else {
		return vs
	}
	iStrSlice := []string{}
	for _, a := range i64c.in {
		iStrSlice = append(iStrSlice, strconv.FormatInt(a, 10))
	}
	return append(vs, fmt.Errorf("%s must be in [%s]", i64c.field, strings.Join(iStrSlice, ", ")))
}
