package govalid

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Validator .
type Validator struct {
	modelStore map[string]*model
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

// Violation checks the struct s against all constraints and custom
// validation functions, if any. It returns an error if the struct
// fails validation. If the type being validated is not a struct,
// ErrNotPtrToStruct will be returned. If the type being validated
// has not yet been registered, ErrNotRegistered is returned.
func (v *Validator) Violation(s interface{}) (string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return "", ErrNotPtrToStruct
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return "", ErrNotPtrToStruct
	}
	m := v.modelStore[t.Name()]
	if m == nil {
		return "", ErrNotRegistered
	}
	return m.violation(s), nil
}

// Violations checks the struct s against all constraints and custom
// validation functions, if any. It returns a slice of errors if the
// struct fails validation. If the type being validated is not a
// struct, ErrNotPtrToStruct alone will be returned. If the type
// being validated has not yet been registered, ErrNotRegistered
// alone is returned.
func (v *Validator) Violations(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return nil, ErrNotPtrToStruct
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil, ErrNotPtrToStruct
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
		err = m.registerStringConstraint(field)
	case reflect.Int:
		err = m.registerIntConstraint(field)
	case reflect.Int64:
		err = m.registerInt64Constraint(field)
	case reflect.Float32:
		err = m.registerFloat32Constraint(field)
	case reflect.Float64:
		err = m.registerFloat64Constraint(field)
	case reflect.Struct:
		if _, ok := field.Type.FieldByName("String"); ok {
			err = m.registerStringConstraint(field)
		} else if _, ok := field.Type.FieldByName("Int64"); ok {
			err = m.registerInt64Constraint(field)
		} else if _, ok := field.Type.FieldByName("Float64"); ok {
			err = m.registerFloat64Constraint(field)
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
		v.registerField(m, typ.Field(i))
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
