package govalid

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type Validator struct {
	customRules map[string]func(v any) error
}

func New() *Validator {
	return &Validator{
		customRules: make(map[string]func(v any) error),
	}
}

var defaultValidator = New()

func Rule(name string, validator func(v any) error) {
	defaultValidator.Rule(name, validator)
}

func Validate(v any) error {
	return defaultValidator.Validate(v)
}

func (v *Validator) Rule(name string, validator func(v any) error) {
	v.customRules[name] = validator
}

func (v *Validator) Validate(val any) error {
	rv := reflect.ValueOf(val)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return errors.New("can not validate nil")
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("can not validate value of kind %s", rv.Kind())
	}
	return v.validateStruct(rv, nil)
}

func (v *Validator) validate(val reflect.Value, rules []string) error {
	switch val.Kind() {
	case reflect.Float32, reflect.Float64:
		return validateFloat(val.Float(), rules, v.customRule)
	case reflect.String:
		return validateString(val.String(), rules, v.customRule)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return validateInt(val.Int(), rules, v.customRule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return validateUint(val.Uint(), rules, v.customRule)
	case reflect.Struct:
		return v.validateStruct(val, rules)
	case reflect.Pointer:
		return v.validatePointer(val, rules)
	case reflect.Slice, reflect.Array:
		return v.validateSlice(val, rules)
	}
	return nil
}

func (v *Validator) validateStruct(rv reflect.Value, rules []string) error {
	for _, rule := range rules {
		if err := v.customRule(rv.Interface(), rule); err != nil {
			return err
		}
	}
	ty := rv.Type()
	for i := 0; i < ty.NumField(); i++ {
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
		if err := v.validate(fv, parts); err != nil {
			return wrap(fmt.Sprintf("field %s", sf.Name), err)
		}
	}
	return nil
}

func (v *Validator) validatePointer(val reflect.Value, rules []string) error {
	req := isReq(rules)
	if req && val.IsNil() {
		return NewValidationError("required")
	}
	if !req && val.IsNil() {
		return nil
	}
	for i, rule := range rules {
		if rule == "dive" {
			if !val.IsZero() && i < len(rules) {
				return v.validate(val.Elem(), rules[i+1:])
			}
			return nil
		}
		if err := v.customRule(val.Interface(), rule); err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) validateSlice(val reflect.Value, rules []string) error {
	req := isReq(rules)
	if req && val.IsNil() {
		return NewValidationError("required")
	}
	if !req && val.IsNil() {
		return nil
	}
	for i, rule := range rules {
		if rule == "dive" {
			if !val.IsZero() {
				for j := 0; j < val.Len(); j++ {
					if err := v.validate(val.Index(j), rules[i+1:]); err != nil {
						return wrap(fmt.Sprintf("index %d", j), err)
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
			if uint64(val.Len()) > max {
				return NewValidationError(fmt.Sprintf("max %d", max))
			}
			continue
		}
		min, ok, err := getUintSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if uint64(val.Len()) < min {
				return NewValidationError(fmt.Sprintf("min %d", min))
			}
			continue
		}
		if err := v.customRule(val.Interface(), rule); err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) customRule(val any, rule string) error {
	if f, ok := v.customRules[rule]; ok {
		return f(val)
	}
	return nil
}

func validateFloat(v float64, rules []string, customRule func(v any, rule string) error) error {
	req := isReq(rules)
	if req && v == 0 {
		return NewValidationError("required")
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
				return NewValidationError(fmt.Sprintf("max %f", max))
			}
			continue
		}
		min, ok, err := getFloatSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if v < min {
				return NewValidationError(fmt.Sprintf("min %f", min))
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateInt(v int64, rules []string, customRule func(v any, rule string) error) error {
	req := isReq(rules)
	if req && v == 0 {
		return NewValidationError("required")
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
				return NewValidationError(fmt.Sprintf("max %d", max))
			}
			continue
		}
		min, ok, err := getIntSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if v < min {
				return NewValidationError(fmt.Sprintf("min %d", min))
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateUint(v uint64, rules []string, customRule func(v any, rule string) error) error {
	req := isReq(rules)
	if req && v == 0 {
		return NewValidationError("required")
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
				return NewValidationError(fmt.Sprintf("max %d", max))
			}
			continue
		}
		min, ok, err := getUintSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if v < min {
				return NewValidationError(fmt.Sprintf("min %d", min))
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateString(v string, rules []string, customRule func(v any, rule string) error) error {
	req := isReq(rules)
	if req && v == "" {
		return NewValidationError("required")
	}
	if !req && v == "" {
		return nil
	}
	for _, rule := range rules {
		if rule == "email" {
			if !strings.Contains(v, "@") {
				return NewValidationError("email")
			}
			continue
		}
		max, ok, err := getIntSize(rule, "max")
		if err != nil {
			return err
		}
		if ok {
			if len(v) > int(max) {
				return NewValidationError(fmt.Sprintf("max %d", max))
			}
			continue
		}
		min, ok, err := getIntSize(rule, "min")
		if err != nil {
			return err
		}
		if ok {
			if len(v) < int(min) {
				return NewValidationError(fmt.Sprintf("min %d", min))
			}
			continue
		}
		if err := customRule(v, rule); err != nil {
			return err
		}
	}
	return nil
}

func isReq(rules []string) bool {
	return slices.Contains(rules, "required")
}

func getIntSize(rule string, prefix string) (int64, bool, error) {
	if !strings.HasPrefix(rule, prefix+":") {
		return 0, false, nil
	}
	v, err := strconv.ParseInt(rule[len(prefix)+1:], 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid %s rule: %w", prefix, err)
	}
	return v, true, nil
}

func getUintSize(rule string, prefix string) (uint64, bool, error) {
	if !strings.HasPrefix(rule, prefix+":") {
		return 0, false, nil
	}
	v, err := strconv.ParseUint(rule[len(prefix)+1:], 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid %s rule: %w", prefix, err)
	}
	return v, true, nil
}

func getFloatSize(rule string, prefix string) (float64, bool, error) {
	if !strings.HasPrefix(rule, prefix+":") {
		return 0, false, nil
	}
	v, err := strconv.ParseFloat(rule[len(prefix)+1:], 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid %s rule: %w", prefix, err)
	}
	return v, true, nil
}
