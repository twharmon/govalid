package govalid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagKey = "govalid"

type model struct {
	name        string
	constraints []constraint
	custom      []func(interface{}) error
}

var modelStore map[string]*model

func (m *model) addToRegistry(name string) {
	if modelStore == nil {
		modelStore = make(map[string]*model)
	}
	if modelStore[name] != nil {
		panic(fmt.Sprintf("%s is already registered", name))
	}
	modelStore[name] = m
}

func (m *model) violation(s interface{}) error {
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i, c := range m.constraints {
		if err := c.violation(val.Field(i)); err != nil {
			return err
		}
	}
	for _, v := range m.custom {
		if err := v(s); err != nil {
			return err
		}
	}
	return nil
}

func (m *model) violations(s interface{}) []error {
	var vs []error
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i, c := range m.constraints {
		vs = append(vs, c.violations(val.Field(i))...)
	}
	for _, v := range m.custom {
		if err := v(s); err != nil {
			vs = append(vs, err)
		}
	}
	return vs
}

func (m *model) registerStringConstraint(field reflect.StructField) {
	c := new(stringConstraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok := m.getIntFromTag(field, tag, "max"); ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok := m.getIntFromTag(field, tag, "min"); ok {
		c.isMinSet = true
		c.min = min
	}
	if in, ok := m.getStringFromTag(tag, "in"); ok {
		c.in = strings.Split(in, ",")
	}
	if reStr, ok := m.getStringFromTag(tag, "regex"); ok {
		re, err := regexp.Compile(reStr)
		if err != nil {
			panic(fmt.Sprintf("govalid model registration error (%s.%s `regex:%s`): %s", m.name, field.Name, reStr, err.Error()))
		}
		c.regex = re
	}
	m.constraints = append(m.constraints, c)
}

func (m *model) registerIntConstraint(field reflect.StructField) {
	c := new(intConstraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok := m.getIntFromTag(field, tag, "max"); ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok := m.getIntFromTag(field, tag, "min"); ok {
		c.isMinSet = true
		c.min = min
	}
	if in, ok := m.getStringFromTag(tag, "in"); ok {
		for _, optStr := range strings.Split(in, ",") {
			opt, err := strconv.Atoi(optStr)
			if err != nil {
				panic(fmt.Sprintf("govalid model registration error (%s.%s `in:%s`): %s", m.name, field.Name, in, err.Error()))
			}
			c.in = append(c.in, opt)
		}
	}
	m.constraints = append(m.constraints, c)
}

func (m *model) registerInt64Constraint(field reflect.StructField) {
	c := new(int64Constraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok := m.getInt64FromTag(field, tag, "max"); ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok := m.getInt64FromTag(field, tag, "min"); ok {
		c.isMinSet = true
		c.min = min
	}
	if in, ok := m.getStringFromTag(tag, "in"); ok {
		for _, optStr := range strings.Split(in, ",") {
			opt, err := strconv.ParseInt(optStr, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("govalid model registration error (%s.%s `in:%s`): %s", m.name, field.Name, in, err.Error()))
			}
			c.in = append(c.in, opt)
		}
	}
	m.constraints = append(m.constraints, c)
}

func (m *model) registerFloat64Constraint(field reflect.StructField) {
	c := new(float64Constraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok := m.getFloat64FromTag(field, tag, "max"); ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok := m.getFloat64FromTag(field, tag, "min"); ok {
		c.isMinSet = true
		c.min = min
	}
	m.constraints = append(m.constraints, c)
}

func (m *model) registerFloat32Constraint(field reflect.StructField) {
	c := new(float32Constraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok := m.getFloat32FromTag(field, tag, "max"); ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok := m.getFloat32FromTag(field, tag, "min"); ok {
		c.isMinSet = true
		c.min = min
	}
	m.constraints = append(m.constraints, c)
}

func (m *model) registerTimeConstraint(field reflect.StructField) {
	c := new(timeConstraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok := m.getInt64FromTag(field, tag, "max"); ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok := m.getInt64FromTag(field, tag, "min"); ok {
		c.isMinSet = true
		c.min = min
	}
	m.constraints = append(m.constraints, c)
}

func (m *model) registerNilConstraint() {
	c := new(nilConstraint)
	m.constraints = append(m.constraints, c)
}

func (m *model) getBoolFromTag(tag string, key string) bool {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		if c == key {
			return true
		}
	}
	return false
}

func (m *model) getIntFromTag(field reflect.StructField, tag string, key string) (int, bool) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			i, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(fmt.Sprintf("govalid model registration error (%s.%s `%s`): %s", m.name, field.Name, c, err.Error()))
			}
			return i, true
		}
	}
	return 0, false
}

func (m *model) getInt64FromTag(field reflect.StructField, tag string, key string) (int64, bool) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			i, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				panic(fmt.Sprintf("govalid model registration error (%s.%s `%s`): %s", m.name, field.Name, c, err.Error()))
			}
			return i, true
		}
	}
	return 0, false
}

func (m *model) getFloat32FromTag(field reflect.StructField, tag string, key string) (float32, bool) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			f, err := strconv.ParseFloat(parts[1], 32)
			if err != nil {
				panic(fmt.Sprintf("govalid model registration error (%s.%s `%s`): %s", m.name, field.Name, c, err.Error()))
			}
			return float32(f), true
		}
	}
	return 0, false
}

func (m *model) getFloat64FromTag(field reflect.StructField, tag string, key string) (float64, bool) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			f, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				panic(fmt.Sprintf("govalid model registration error (%s.%s `%s`): %s", m.name, field.Name, c, err.Error()))
			}
			return f, true
		}
	}
	return 0, false
}

func (m *model) getStringFromTag(tag string, key string) (string, bool) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			return parts[1], true
		}
	}
	return "", false
}
