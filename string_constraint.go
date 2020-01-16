package govalid

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
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

func (sc *stringConstraint) error(val reflect.Value) error {
	empty := true
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
		return fmt.Errorf("%w: %s is required", ErrInvalidStruct, sc.field)
	}
	strLen := utf8.RuneCountInString(s)
	if sc.isMaxSet && strLen > sc.max {
		return fmt.Errorf("%w: %s can not be longer than %d characters", ErrInvalidStruct, sc.field, sc.max)
	}
	if sc.isMinSet && strLen < sc.min {
		return fmt.Errorf("%w: %s must be at least %d characters", ErrInvalidStruct, sc.field, sc.min)
	}
	if sc.regex != nil && !sc.regex.MatchString(s) {
		return fmt.Errorf("%w: %s must match regex /%s/", ErrInvalidStruct, sc.field, sc.regex.String())
	}
	if len(sc.in) > 0 {
		for _, opt := range sc.in {
			if s == opt {
				return nil
			}
		}
	} else {
		return nil
	}
	return fmt.Errorf("%w: %s must be in [%s]", ErrInvalidStruct, sc.field, strings.Join(sc.in, ", "))
}

func (sc *stringConstraint) errors(val reflect.Value) []error {
	var vs []error
	empty := true
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
		vs = append(vs, fmt.Errorf("%w: %s is required", ErrInvalidStruct, sc.field))
	}
	strLen := utf8.RuneCountInString(s)
	if sc.isMaxSet && strLen > sc.max {
		vs = append(vs, fmt.Errorf("%w: %s can not be longer than %d characters", ErrInvalidStruct, sc.field, sc.max))
	}
	if sc.isMinSet && strLen < sc.min {
		vs = append(vs, fmt.Errorf("%w: %s must be at least %d characters", ErrInvalidStruct, sc.field, sc.min))
	}
	if sc.regex != nil && !sc.regex.MatchString(s) {
		vs = append(vs, fmt.Errorf("%w: %s must match regex /%s/", ErrInvalidStruct, sc.field, sc.regex.String()))
	}
	if len(sc.in) > 0 {
		for _, opt := range sc.in {
			if s == opt {
				return vs
			}
		}
	} else {
		return vs
	}
	return append(vs, fmt.Errorf("%w: %s must be in [%s]", ErrInvalidStruct, sc.field, strings.Join(sc.in, ", ")))
}
