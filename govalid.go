package govalid

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var customRules = make(map[string]func(v any) error)

func Rule(name string, validator func(v any) error) {
	customRules[name] = validator
}

func Validate(v any) error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return errors.New("can not validate nil")
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("can not validate value of kind %s", rv.Kind())
	}
	return validateStruct(rv, nil)
}

func validate(v reflect.Value, rules []string) error {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return validateFloat(v.Float(), rules)
	case reflect.String:
		return validateString(v.String(), rules)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return validateInt(v.Int(), rules)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return validateUint(v.Uint(), rules)
	case reflect.Struct:
		return validateStruct(v, rules)
	case reflect.Pointer:
		return validatePointer(v, rules)
	case reflect.Slice, reflect.Array:
		return validateSlice(v, rules)
	}
	return nil
}

func validateStruct(rv reflect.Value, rules []string) error {
	for _, rule := range rules {
		if err := customRule(rv.Interface(), rule); err != nil {
			return err
		}
	}
	ty := rv.Type()
	for i := range ty.NumField() {
		sf := ty.Field(i)
		if !sf.IsExported() {
			continue
		}
		st := sf.Tag
		tag, ok := st.Lookup("valid")
		if !ok {
			continue
		}
		fv := rv.Field(i)
		parts := strings.Split(tag, "|")
		if err := validate(fv, parts); err != nil {
			return fmt.Errorf("field %s: %w", sf.Name, err)
		}
	}
	return nil
}

func validatePointer(v reflect.Value, rules []string) error {
	req := isReq(rules)
	if req && v.IsNil() {
		return errors.New("required")
	}
	if !req && v.IsNil() {
		return nil
	}
	for i, rule := range rules {
		if rule == "dive" {
			if !v.IsZero() && i < len(rules) {
				return validate(v.Elem(), rules[i+1:])
			}
			return nil
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateSlice(v reflect.Value, rules []string) error {
	req := isReq(rules)
	if req && v.IsNil() {
		return errors.New("required")
	}
	if !req && v.IsNil() {
		return nil
	}
	for i, rule := range rules {
		if rule == "dive" {
			if !v.IsZero() && i < len(rules) {
				for j := range v.Len() {
					if err := validate(v.Index(j), rules[i+1:]); err != nil {
						return fmt.Errorf("index %d: %w", j, err)
					}
				}
			}
			return nil
		}
		max, ok, err := getUintSize(rule, "max")
		if err != nil {
			return err
		}
		if ok {
			if uint64(v.Len()) > max {
				return fmt.Errorf("max %d", max)
			}
			continue
		}
		min, ok, err := getUintSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if uint64(v.Len()) < min {
				return fmt.Errorf("min %d", min)
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateFloat(v float64, rules []string) error {
	req := isReq(rules)
	if req && v == 0 {
		return errors.New("required")
	}
	if !req && v == 0 {
		return nil
	}
	for _, rule := range rules {
		max, ok, err := getFloatSize(rule, "max")
		if err != nil {
			return err
		}
		if ok {
			if v > max {
				return fmt.Errorf("max %f", max)
			}
			continue
		}
		min, ok, err := getFloatSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if v < min {
				return fmt.Errorf("min %f", min)
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateInt(v int64, rules []string) error {
	req := isReq(rules)
	if req && v == 0 {
		return errors.New("required")
	}
	if !req && v == 0 {
		return nil
	}
	for _, rule := range rules {
		max, ok, err := getIntSize(rule, "max")
		if err != nil {
			return err
		}
		if ok {
			if v > max {
				return fmt.Errorf("max %d", max)
			}
			continue
		}
		min, ok, err := getIntSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if v < min {
				return fmt.Errorf("min %d", min)
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateUint(v uint64, rules []string) error {
	req := isReq(rules)
	if req && v == 0 {
		return errors.New("required")
	}
	if !req && v == 0 {
		return nil
	}
	for _, rule := range rules {
		max, ok, err := getUintSize(rule, "max")
		if err != nil {
			return err
		}
		if ok {
			if v > max {
				return fmt.Errorf("max %d", max)
			}
			continue
		}
		min, ok, err := getUintSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if v < min {
				return fmt.Errorf("min %d", min)
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateString(v string, rules []string) error {
	req := isReq(rules)
	if req && v == "" {
		return errors.New("required")
	}
	if !req && v == "" {
		return nil
	}
	for _, rule := range rules {
		max, ok, err := getUintSize(rule, "max")
		if err != nil {
			return err
		}
		if ok {
			if uint64(len(v)) > max {
				return fmt.Errorf("max %d", max)
			}
			continue
		}
		min, ok, err := getUintSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if uint64(len(v)) < min {
				return fmt.Errorf("min %d", min)
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func customRule(v any, rule string) error {
	if validator, ok := customRules[rule]; ok {
		if err := validator(v); err != nil {
			return err
		}
	}
	return nil
}

func getIntSize(rule string, ty string) (int64, bool, error) {
	prefix := fmt.Sprintf("%s:", ty)
	if strings.HasPrefix(rule, prefix) {
		s := strings.TrimPrefix(rule, prefix)
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, false, err
		}
		return i, true, nil
	}
	return 0, false, nil
}

func getUintSize(rule string, ty string) (uint64, bool, error) {
	prefix := fmt.Sprintf("%s:", ty)
	if strings.HasPrefix(rule, prefix) {
		s := strings.TrimPrefix(rule, prefix)
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, false, err
		}
		return i, true, nil
	}
	return 0, false, nil
}

func getFloatSize(rule string, ty string) (float64, bool, error) {
	prefix := fmt.Sprintf("%s:", ty)
	if strings.HasPrefix(rule, prefix) {
		s := strings.TrimPrefix(rule, prefix)
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, false, err
		}
		return i, true, nil
	}
	return 0, false, nil
}

func isReq(rules []string) bool {
	for _, rule := range rules {
		if rule == "req" {
			return true
		}
	}
	return false
}
