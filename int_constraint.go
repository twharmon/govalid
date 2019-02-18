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

func makeIntConstraint(name string, field reflect.StructField) {
	ic := new(intConstraint)
	ic.field = strings.ToLower(field.Name)
	req, ok := field.Tag.Lookup("req")
	if ok {
		ic.req = req == "true"
	}
	maxStr, ok := field.Tag.Lookup("max")
	if ok {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			panic(err)
		}
		ic.max = max
	}
	minStr, ok := field.Tag.Lookup("min")
	if ok {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			panic(err)
		}
		ic.min = min
	}
	inStr, ok := field.Tag.Lookup("in")
	if ok {
		inStrSlice := strings.Split(inStr, ",")
		inIntSlice := []int{}
		for _, iStr := range inStrSlice {
			i, err := strconv.Atoi(iStr)
			if err != nil {
				panic(err)
			}
			inIntSlice = append(inIntSlice, i)
		}
		ic.in = inIntSlice
	}
	modelStore.add(name, ic)
}
