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

func (ic *intConstraint) validate(val reflect.Value) string {
	i := val.Interface().(int)
	if !ic.req && i == 0 {
		return ""
	}
	if ic.req && i == 0 {
		return fmt.Sprintf("%s is required", ic.field)
	}
	if ic.max > 0 && i > ic.max {
		return fmt.Sprintf("%s can not be greater than %d", ic.field, ic.max)
	}
	if ic.min > 0 && i < ic.min {
		return fmt.Sprintf("%s must be at least %d", ic.field, ic.min)
	}
	if len(ic.in) > 0 {
		for _, opt := range ic.in {
			if i == opt {
				return ""
			}
		}
	} else {
		return ""
	}
	iStrSlice := []string{}
	for _, a := range ic.in {
		iStrSlice = append(iStrSlice, strconv.Itoa(a))
	}
	return fmt.Sprintf("%s must be in [%s]", ic.field, strings.Join(iStrSlice, ", "))
}
