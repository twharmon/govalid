package govalid

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type int64Constraint struct {
	field string
	req   bool
	min   int64
	max   int64
	in    []int64
}

func (i64c *int64Constraint) validate(val reflect.Value) []string {
	i64 := val.Interface().(int64)
	if !i64c.req && i64 == 0 {
		return nil
	}
	if i64c.req && i64 == 0 {
		return []string{fmt.Sprintf("%s is required", i64c.field)}
	}
	var vs []string
	if i64c.max > 0 && i64 > i64c.max {
		vs = append(vs, fmt.Sprintf("%s can not be greater than %d", i64c.field, i64c.max))
	}
	if i64c.min > 0 && i64 < i64c.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %d", i64c.field, i64c.min))
	}
	if len(i64c.in) > 0 {
		in := false
		for _, opt := range i64c.in {
			if i64 == opt {
				in = true
				break
			}
		}
		if !in {
			iStrSlice := []string{}
			for _, a := range i64c.in {
				iStrSlice = append(iStrSlice, strconv.FormatInt(a, 10))
			}
			vs = append(vs, fmt.Sprintf("%s must be in [%s]", i64c.field, strings.Join(iStrSlice, ", ")))
		}
	}
	return vs
}
