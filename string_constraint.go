package govalid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type stringConstraint struct {
	field string
	req   bool
	min   int
	max   int
	in    []string
	regex *regexp.Regexp
}

func (sc *stringConstraint) validate(val reflect.Value) []string {
	s := val.Interface().(string)
	var vs []string
	if !sc.req && s == "" {
		return vs
	}
	if sc.req && s == "" {
		vs = append(vs, fmt.Sprintf("%s is required", sc.field))
	}
	strLen := utf8.RuneCountInString(s)
	if sc.max > 0 && strLen > sc.max {
		vs = append(vs, fmt.Sprintf("%s can not be longer than %d characters", sc.field, sc.max))
	}
	if sc.min > 0 && strLen < sc.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %d characters", sc.field, sc.min))
	}
	if sc.regex != nil && !sc.regex.MatchString(s) {
		vs = append(vs, fmt.Sprintf("%s must match regex /%s/", sc.field, sc.regex.String()))
	}
	if len(sc.in) > 0 {
		in := false
		for _, opt := range sc.in {
			if s == opt {
				in = true
				break
			}
		}
		if !in {
			vs = append(vs, fmt.Sprintf("%s must be in [%s]", sc.field, strings.TrimSuffix(strings.Join(sc.in, ", "), ", ")))
		}
	}
	return vs
}

func makeStringConstraint(name string, field reflect.StructField) {
	sc := new(stringConstraint)
	sc.field = strings.ToLower(field.Name)
	req, ok := field.Tag.Lookup("req")
	if ok {
		sc.req = req == "true"
	}
	maxStr, ok := field.Tag.Lookup("max")
	if ok {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			panic(err)
		}
		sc.max = max
	}
	minStr, ok := field.Tag.Lookup("min")
	if ok {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			panic(err)
		}
		sc.min = min
	}
	inStr, ok := field.Tag.Lookup("in")
	if ok {
		in := strings.Split(inStr, ",")
		sc.in = in
	}
	regex, ok := field.Tag.Lookup("regex")
	if ok {
		re, err := regexp.Compile(regex)
		if err != nil {
			panic(err)
		}
		sc.regex = re
	}
	modelStore.add(name, sc)
}
