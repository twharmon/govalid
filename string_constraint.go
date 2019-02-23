package govalid

import (
	"fmt"
	"reflect"
	"regexp"
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

func (sc *stringConstraint) validate(val reflect.Value) string {
	s := val.Interface().(string)
	if !sc.req && s == "" {
		return ""
	}
	if sc.req && s == "" {
		return fmt.Sprintf("%s is required", sc.field)
	}
	strLen := utf8.RuneCountInString(s)
	if sc.max > 0 && strLen > sc.max {
		return fmt.Sprintf("%s can not be longer than %d characters", sc.field, sc.max)
	}
	if sc.min > 0 && strLen < sc.min {
		return fmt.Sprintf("%s must be at least %d characters", sc.field, sc.min)
	}
	if sc.regex != nil && !sc.regex.MatchString(s) {
		return fmt.Sprintf("%s must match regex /%s/", sc.field, sc.regex.String())
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
			return fmt.Sprintf("%s must be in [%s]", sc.field, strings.Join(sc.in, ", "))
		}
	}
	return ""
}
