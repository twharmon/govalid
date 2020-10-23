package govalid

import (
	"fmt"
	"reflect"
)

type int64Constraint struct {
	field       string
	req         bool
	isMinSet    bool
	min         int64
	isMaxSet    bool
	max         int64
	in          []int64
	customRules []func(string, int64) string
}

func (i64c *int64Constraint) violation(val reflect.Value) string {
	var empty bool
	i64, ok := val.Interface().(int64)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		i64 = val.FieldByName("Int64").Interface().(int64)
		empty = false
	} else {
		empty = i64 == 0
	}
	if !i64c.req && empty {
		return ""
	}
	if i64c.req && empty {
		return fmt.Sprintf("%s is required", i64c.field)
	}
	if i64c.isMaxSet && i64 > i64c.max {
		return fmt.Sprintf("%s can not be greater than %d", i64c.field, i64c.max)
	}
	if i64c.isMinSet && i64 < i64c.min {
		return fmt.Sprintf("%s must be at least %d", i64c.field, i64c.min)
	}
	if !containsInt64(i64c.in, i64) {
		return fmt.Sprintf("%s must be in %v", i64c.field, i64c.in)
	}
	for _, f := range i64c.customRules {
		if vio := f(i64c.field, i64); vio != "" {
			return vio
		}
	}
	return ""
}

func (i64c *int64Constraint) violations(val reflect.Value) []string {
	var vs []string
	var empty bool
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
		vs = append(vs, fmt.Sprintf("%s is required", i64c.field))
	}
	if i64c.isMaxSet && i64 > i64c.max {
		vs = append(vs, fmt.Sprintf("%s can not be greater than %d", i64c.field, i64c.max))
	}
	if i64c.isMinSet && i64 < i64c.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %d", i64c.field, i64c.min))
	}
	if !containsInt64(i64c.in, i64) {
		vs = append(vs, fmt.Sprintf("%s must be in %v", i64c.field, i64c.in))
	}
	for _, f := range i64c.customRules {
		if vio := f(i64c.field, i64); vio != "" {
			vs = append(vs, vio)
		}
	}
	return vs
}

func containsInt64(a []int64, x int64) bool {
	if len(a) == 0 {
		return true
	}

	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
