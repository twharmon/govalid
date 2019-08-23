package govalid

import (
	"fmt"
	"reflect"
	"time"
)

type timeConstraint struct {
	field    string
	req      bool
	isMinSet bool
	min      int64
	isMaxSet bool
	max      int64
}

func (tc *timeConstraint) violation(val reflect.Value) error {
	empty := true
	var tUnix int64
	t, ok := val.Interface().(time.Time)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		t = val.FieldByName("Time").Interface().(time.Time)
		empty = false
		tUnix = t.Unix()
	} else {
		tUnix = t.Unix()
		empty = tUnix < 0
	}
	if !tc.req && empty {
		return nil
	}
	if tc.req && empty {
		return fmt.Errorf("%s is required", tc.field)
	}
	age := time.Now().Unix() - tUnix
	if tc.isMaxSet && age > tc.max {
		return fmt.Errorf("%s can not have age greater than %d seconds", tc.field, tc.max)
	}
	if tc.isMinSet && age < tc.min {
		return fmt.Errorf("%s must have age at least %d seconds", tc.field, tc.min)
	}
	return nil
}

func (tc *timeConstraint) violations(val reflect.Value) []error {
	var vs []error
	empty := true
	var tUnix int64
	t, ok := val.Interface().(time.Time)
	if !ok && val.FieldByName("Valid").Interface().(bool) {
		t = val.FieldByName("Time").Interface().(time.Time)
		empty = false
		tUnix = t.Unix()
	} else {
		tUnix = t.Unix()
		empty = tUnix < 0
	}
	if !tc.req && empty {
		return nil
	}
	if tc.req && empty {
		vs = append(vs, fmt.Errorf("%s is required", tc.field))
	}
	age := time.Now().Unix() - tUnix
	if tc.isMaxSet && age > tc.max {
		vs = append(vs, fmt.Errorf("%s can not have age greater than %d seconds", tc.field, tc.max))
	}
	if tc.isMinSet && age < tc.min {
		vs = append(vs, fmt.Errorf("%s must have age at least %d seconds", tc.field, tc.min))
	}
	return vs
}
