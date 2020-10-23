package govalid

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Validator .
type Validator struct {
	modelStore   map[string]*model
	stringRules  map[string]func(string, string) string
	int64Rules   map[string]func(string, int64) string
	float64Rules map[string]func(string, float64) string
}

// Register is required for all structs that you wish
// to validate. It is intended to be ran at load time
// and caches information about the structs to reduce
// run time allocations.
func (v *Validator) Register(structs ...interface{}) error {
	for _, s := range structs {
		if err := v.register(s); err != nil {
			return err
		}
	}
	return nil
}

// AddCustom adds custom validation functions to struct s.
func (v *Validator) AddCustom(s interface{}, f ...func(interface{}) string) error {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	n := t.Name()
	m := v.modelStore[n]
	if m == nil {
		return ErrNotRegistered
	}
	m.custom = append(m.custom, f...)
	return nil
}

// AddCustomStringRule adds custom validation tag for string.
func (v *Validator) AddCustomStringRule(name string, validatorFunc func(string, string) string) {
	v.stringRules[name] = validatorFunc
}

// AddCustomInt64Rule adds custom validation tag for int64.
func (v *Validator) AddCustomInt64Rule(name string, validatorFunc func(string, int64) string) {
	v.int64Rules[name] = validatorFunc
}

// AddCustomFloat64Rule adds custom validation tag for float64.
func (v *Validator) AddCustomFloat64Rule(name string, validatorFunc func(string, float64) string) {
	v.float64Rules[name] = validatorFunc
}

// Violation checks the struct s against all constraints and custom
// validation functions, if any. It returns an violation if the
// struct fails validation. If the type being validated is not a
// struct, ErrNotStruct will be returned. If the type being validated
// has not yet been registered, ErrNotRegistered is returned.
func (v *Validator) Violation(s interface{}) (string, error) {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", ErrNotStruct
	}
	m := v.modelStore[t.Name()]
	if m == nil {
		return "", ErrNotRegistered
	}
	return m.violation(s), nil
}

// Violations checks the struct s against all constraints and custom
// validation functions, if any. It returns a slice of violations if
// the struct fails validation. If the type being validated is not a
// struct, ErrNotStruct will be returned. If the type being validated
// has not yet been registered, ErrNotRegistered is returned.
func (v *Validator) Violations(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	m := v.modelStore[t.Name()]
	if m == nil {
		return nil, ErrNotRegistered
	}
	return m.violations(s), nil
}

func (v *Validator) registerField(m *model, field reflect.StructField) error {
	firstLetter := string(field.Name[0])
	if firstLetter != strings.ToUpper(firstLetter) {
		m.registerNilConstraint()
		return nil
	}
	var err error
	switch field.Type.Kind() {
	case reflect.String:
		err = m.registerStringConstraint(field, v.stringRules)
	case reflect.Int:
		err = m.registerIntConstraint(field)
	case reflect.Int64:
		err = m.registerInt64Constraint(field, v.int64Rules)
	case reflect.Float32:
		err = m.registerFloat32Constraint(field)
	case reflect.Float64:
		err = m.registerFloat64Constraint(field, v.float64Rules)
	case reflect.Struct:
		if _, ok := field.Type.FieldByName("String"); ok {
			err = m.registerStringConstraint(field, v.stringRules)
		} else if _, ok := field.Type.FieldByName("Int64"); ok {
			err = m.registerInt64Constraint(field, v.int64Rules)
		} else if _, ok := field.Type.FieldByName("Float64"); ok {
			err = m.registerFloat64Constraint(field, v.float64Rules)
		} else {
			m.registerNilConstraint()
		}
	default:
		m.registerNilConstraint()
	}
	return err
}

func (v *Validator) register(s interface{}) error {
	typ := reflect.TypeOf(s)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return errors.New("only structs can be registered")
	}
	name := typ.Name()
	m := new(model)
	m.name = name
	for i := 0; i < typ.NumField(); i++ {
		if err := v.registerField(m, typ.Field(i)); err != nil {
			return err
		}
	}
	return v.addModelToRegistry(m, name)
}

func (v *Validator) addModelToRegistry(m *model, name string) error {
	if v.modelStore[name] != nil {
		return fmt.Errorf("%s is already registered", name)
	}
	v.modelStore[name] = m
	return nil
}
