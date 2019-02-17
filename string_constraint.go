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

func (sc *stringConstraint) validate(val reflect.Value) error {
	s := val.Interface().(string)
	if !sc.req && s == "" {
		return nil
	}
	if sc.req && s == "" {
		return fmt.Errorf("%s is required", sc.field)
	}
	strLen := utf8.RuneCountInString(s)
	if sc.max > 0 && strLen > sc.max {
		return fmt.Errorf("%s can not be longer than %d characters", sc.field, sc.max)
	}
	if sc.min > 0 && strLen < sc.min {
		return fmt.Errorf("%s must be at least %d characters", sc.field, sc.min)
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
			return fmt.Errorf("%s must be in [%s]", sc.field, strings.TrimSuffix(strings.Join(sc.in, ", "), ", "))
		}
	}
	if sc.regex != nil && !sc.regex.MatchString(s) {
		return fmt.Errorf("%s must match regex /%s/", sc.field, sc.regex.String())
	}
	return nil
}

func makeStringConstraint(name string, field reflect.StructField) error {
	sc := new(stringConstraint)
	sc.field = field.Name
	req, ok := field.Tag.Lookup("req")
	if ok {
		sc.req = req == "true"
	}
	maxStr, ok := field.Tag.Lookup("max")
	if ok {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return err
		}
		sc.max = max
	}
	minStr, ok := field.Tag.Lookup("min")
	if ok {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return err
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
			return err
		}
		sc.regex = re
	}
	store.Add(name, sc)
	return nil
}
