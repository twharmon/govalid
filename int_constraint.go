package govalid

import (
	"fmt"
	"reflect"
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

func (ic *intConstraint) violation(val reflect.Value) string {
	i := val.Interface().(int)
	if !ic.req && i == 0 {
		return ""
	}
	if ic.req && i == 0 {
		return fmt.Sprintf("%s is required", ic.field)
	}
	if ic.isMaxSet && i > ic.max {
		return fmt.Sprintf("%s can not be greater than %d", ic.field, ic.max)
	}
	if ic.isMinSet && i < ic.min {
		return fmt.Sprintf("%s must be at least %d", ic.field, ic.min)
	}
	if !containsInt(ic.in, i) {
		return fmt.Sprintf("%s must be in %v", ic.field, ic.in)
	}

	return ""
}

func (ic *intConstraint) violations(val reflect.Value) []string {
	var vs []string
	i := val.Interface().(int)
	if !ic.req && i == 0 {
		return nil
	}
	if ic.req && i == 0 {
		vs = append(vs, fmt.Sprintf("%s is required", ic.field))
	}
	if ic.isMaxSet && i > ic.max {
		vs = append(vs, fmt.Sprintf("%s can not be greater than %d", ic.field, ic.max))
	}
	if ic.isMinSet && i < ic.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %d", ic.field, ic.min))
	}
	if !containsInt(ic.in, i) {
		vs = append(vs, fmt.Sprintf("%s must be in %v", ic.field, ic.in))
	}

	return vs
}

func containsInt(a []int, x int) bool {
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
