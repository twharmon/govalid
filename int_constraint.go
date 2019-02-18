package govalid

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type intConstraint struct {
	field string
	req   bool
	min   int
	max   int
	in    []int
}

func (ic *intConstraint) validate(val reflect.Value) []string {
	i := val.Interface().(int)
	var vs []string
	if !ic.req && i == 0 {
		return vs
	}
	if ic.req && i == 0 {
		vs = append(vs, fmt.Sprintf("%s is required", ic.field))
		return vs
	}
	if ic.max > 0 && i > ic.max {
		vs = append(vs, fmt.Sprintf("%s can not be greater than %d", ic.field, ic.max))
	}
	if ic.min > 0 && i < ic.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %d", ic.field, ic.min))
	}
	if len(ic.in) > 0 {
		in := false
		for _, opt := range ic.in {
			if i == opt {
				in = true
				break
			}
		}
		if !in {
			iStrSlice := []string{}
			for _, a := range ic.in {
				iStrSlice = append(iStrSlice, strconv.Itoa(a))
			}
			vs = append(vs, fmt.Sprintf("%s must be in [%s]", ic.field, strings.TrimSuffix(strings.Join(iStrSlice, ", "), ", ")))
		}
	}
	return vs
}
