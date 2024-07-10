package govalid

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var customRules = make(map[string]func(v any) (string, error))

func Rule(name string, validator func(v any) (string, error)) {
	customRules[name] = validator
}

func Validate(v any) (string, error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return "", fmt.Errorf("can not validate nil")
	}
	if rv.Kind() != reflect.Struct {
		return "", fmt.Errorf("can not validate value of kind %s", rv.Kind())
	}
	return validateStruct(rv, nil)
}

func validate(v reflect.Value, rules []string) (string, error) {
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
	return "", nil
}

func validateStruct(rv reflect.Value, rules []string) (string, error) {
	for _, rule := range rules {
		if vio, err := customRule(rv.Interface(), rule); err != nil {
			return "", err
		} else if vio != "" {
			return vio, nil
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
		vio, err := validate(fv, parts)
		if err != nil {
			return "", fmt.Errorf("field %s: %w", sf.Name, err)
		}
		if vio != "" {
			return fmt.Sprintf("field: %s: %s", sf.Name, vio), nil
		}
	}
	return "", nil
}

func validatePointer(v reflect.Value, rules []string) (string, error) {
	req := isReq(rules)
	if req && v.IsNil() {
		return "required", nil
	}
	if !req && v.IsNil() {
		return "", nil
	}
	for i, rule := range rules {
		if rule == "dive" {
			if !v.IsZero() && i < len(rules) {
				return validate(v.Elem(), rules[i+1:])
			}
			return "", nil
		}
		if vio, err := customRule(v, rule); err != nil {
			return "", err
		} else if vio != "" {
			return vio, nil
		}
	}
	return "", nil
}

func validateSlice(v reflect.Value, rules []string) (string, error) {
	req := isReq(rules)
	if req && v.IsNil() {
		return "required", nil
	}
	if !req && v.IsNil() {
		return "", nil
	}
	for i, rule := range rules {
		if rule == "dive" {
			if !v.IsZero() && i < len(rules) {
				for i := range v.Len() {
					vio, err := validate(v.Index(i), rules[i+1:])
					if err != nil {
						return "", fmt.Errorf("index %d: %w", i, err)
					}
					if vio != "" {
						return fmt.Sprintf("index %d: %s", i, vio), nil
					}
				}
			}
			return "", nil
		}
		max, ok, err := getUintSize(rule, "max")
		if err != nil {
			return "", err
		}
		if ok {
			if uint64(v.Len()) > max {
				return fmt.Sprintf("max %d", max), nil
			}
			continue
		}
		min, ok, err := getUintSize(rule, "min")
		if err != nil {
			return "", err
		}
		if ok {
			if uint64(v.Len()) < min {
				return fmt.Sprintf("min %d", min), nil
			}
			continue
		}
		if vio, err := customRule(v, rule); err != nil {
			return "", err
		} else if vio != "" {
			return vio, nil
		}
	}
	return "", nil
}

func validateFloat(v float64, rules []string) (string, error) {
	req := isReq(rules)
	if req && v == 0 {
		return "required", nil
	}
	if !req && v == 0 {
		return "", nil
	}
	for _, rule := range rules {
		max, ok, err := getFloatSize(rule, "max")
		if err != nil {
			return "", err
		}
		if ok {
			if v > max {
				return fmt.Sprintf("max %f", max), nil
			}
			continue
		}
		min, ok, err := getFloatSize(rule, "min")
		if err != nil {
			return "", err
		}
		if ok {
			if v < min {
				return fmt.Sprintf("min %f", min), nil
			}
			continue
		}
		if vio, err := customRule(v, rule); err != nil {
			return "", err
		} else if vio != "" {
			return vio, nil
		}
	}
	return "", nil
}

func validateInt(v int64, rules []string) (string, error) {
	req := isReq(rules)
	if req && v == 0 {
		return "required", nil
	}
	if !req && v == 0 {
		return "", nil
	}
	for _, rule := range rules {
		max, ok, err := getIntSize(rule, "max")
		if err != nil {
			return "", err
		}
		if ok {
			if v > max {
				return fmt.Sprintf("max %d", max), nil
			}
			continue
		}
		min, ok, err := getIntSize(rule, "min")
		if err != nil {
			return "", err
		}
		if ok {
			if v < min {
				return fmt.Sprintf("min %d", min), nil
			}
			continue
		}
		if vio, err := customRule(v, rule); err != nil {
			return "", err
		} else if vio != "" {
			return vio, nil
		}
	}
	return "", nil
}

func validateUint(v uint64, rules []string) (string, error) {
	req := isReq(rules)
	if req && v == 0 {
		return "required", nil
	}
	if !req && v == 0 {
		return "", nil
	}
	for _, rule := range rules {
		max, ok, err := getUintSize(rule, "max")
		if err != nil {
			return "", err
		}
		if ok {
			if v > max {
				return fmt.Sprintf("max %d", max), nil
			}
			continue
		}
		min, ok, err := getUintSize(rule, "min")
		if err != nil {
			return "", err
		}
		if ok {
			if v < min {
				return fmt.Sprintf("min %d", min), nil
			}
			continue
		}
		if vio, err := customRule(v, rule); err != nil {
			return "", err
		} else if vio != "" {
			return vio, nil
		}
	}
	return "", nil
}

func validateString(v string, rules []string) (string, error) {
	req := isReq(rules)
	if req && v == "" {
		return "required", nil
	}
	if !req && v == "" {
		return "", nil
	}
	for _, rule := range rules {
		max, ok, err := getUintSize(rule, "max")
		if err != nil {
			return "", err
		}
		if ok {
			if uint64(len(v)) > max {
				return fmt.Sprintf("max %d", max), nil
			}
			continue
		}
		min, ok, err := getUintSize(rule, "min")
		if err != nil {
			return "", err
		}
		if ok {
			if uint64(len(v)) < min {
				return fmt.Sprintf("min %d", min), nil
			}
			continue
		}
		if vio, err := customRule(v, rule); err != nil {
			return "", err
		} else if vio != "" {
			return vio, nil
		}
	}
	return "", nil
}

func customRule(v any, rule string) (string, error) {
	if validator, ok := customRules[rule]; ok {
		vio, err := validator(v)
		if err != nil {
			return "", err
		}
		if vio != "" {
			return vio, nil
		}
	}
	return "", nil
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
