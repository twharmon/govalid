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

func (i64c *int64Constraint) validate(val reflect.Value) error {
	i64 := val.Interface().(int64)
	if !i64c.req && i64 == 0 {
		return nil
	}
	if i64c.req && i64 == 0 {
		return fmt.Errorf("%s is required", i64c.field)
	}
	if i64c.max > 0 && i64 > i64c.max {
		return fmt.Errorf("%s can not be greater than %d", i64c.field, i64c.max)
	}
	if i64c.min > 0 && i64 < i64c.min {
		return fmt.Errorf("%s must be at least %d", i64c.field, i64c.min)
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
			return fmt.Errorf("%s must be in [%s]", i64c.field, strings.TrimSuffix(strings.Join(iStrSlice, ", "), ", "))
		}
	}
	return nil
}

func makeInt64Constraint(name string, field reflect.StructField) error {
	i64c := new(int64Constraint)
	i64c.field = field.Name
	req, ok := field.Tag.Lookup("req")
	if ok {
		i64c.req = req == "true"
	}
	maxStr, ok := field.Tag.Lookup("max")
	if ok {
		max, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			return err
		}
		i64c.max = max
	}
	minStr, ok := field.Tag.Lookup("min")
	if ok {
		min, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return err
		}
		i64c.min = min
	}
	inStr, ok := field.Tag.Lookup("in")
	if ok {
		inStrSlice := strings.Split(inStr, ",")
		inInt64Slice := []int64{}
		for _, iStr := range inStrSlice {
			i64, err := strconv.ParseInt(iStr, 10, 64)
			if err != nil {
				return err
			}
			inInt64Slice = append(inInt64Slice, i64)
		}
		i64c.in = inInt64Slice
	}
	constraintStore.Add(name, i64c)
	return nil
}
