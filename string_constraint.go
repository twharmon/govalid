package govalid

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode/utf8"
)

type stringConstraint struct {
	field    string
	req      bool
	isMinSet bool
	min      int
	isMaxSet bool
	max      int
	in       []string
	regex    *regexp.Regexp
}

func (sc *stringConstraint) violation(val reflect.Value) string {
	var empty bool
	s, ok := val.Interface().(string)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		s = val.FieldByName("String").Interface().(string)
		empty = false
	} else {
		empty = s == ""
	}
	if !sc.req && empty {
		return ""
	}
	if sc.req && empty {
		return fmt.Sprintf("%s is required", sc.field)
	}
	strLen := utf8.RuneCountInString(s)
	if sc.isMaxSet && strLen > sc.max {
		return fmt.Sprintf("%s can not be longer than %d characters", sc.field, sc.max)
	}
	if sc.isMinSet && strLen < sc.min {
		return fmt.Sprintf("%s must be at least %d characters", sc.field, sc.min)
	}
	if sc.regex != nil && !sc.regex.MatchString(s) {
		return fmt.Sprintf("%s must match regex /%s/", sc.field, sc.regex.String())
	}
	return sc.getInViolation(s)
}

func (sc *stringConstraint) violations(val reflect.Value) []string {
	var vs []string
	var empty bool
	s, ok := val.Interface().(string)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		s = val.FieldByName("String").Interface().(string)
		empty = false
	} else {
		empty = s == ""
	}
	if !sc.req && empty {
		return nil
	}
	if sc.req && empty {
		vs = append(vs, fmt.Sprintf("%s is required", sc.field))
	}
	strLen := utf8.RuneCountInString(s)
	if sc.isMaxSet && strLen > sc.max {
		vs = append(vs, fmt.Sprintf("%s can not be longer than %d characters", sc.field, sc.max))
	}
	if sc.isMinSet && strLen < sc.min {
		vs = append(vs, fmt.Sprintf("%s must be at least %d characters", sc.field, sc.min))
	}
	if sc.regex != nil && !sc.regex.MatchString(s) {
		vs = append(vs, fmt.Sprintf("%s must match regex /%s/", sc.field, sc.regex.String()))
	}
	if v := sc.getInViolation(s); v != "" {
		vs = append(vs, v)
	}
	return vs
}

func (sc *stringConstraint) getInViolation(s string) string {
	if len(sc.in) > 0 {
		for _, opt := range sc.in {
			if s == opt {
				return ""
			}
		}
		return fmt.Sprintf("%s must be in %v", sc.field, sc.in)
	}
	return ""
}
