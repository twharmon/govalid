package govalid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagKey = "validate"

type model struct {
	name        string // todo: .
	constraints []constraint
	custom      []func(interface{}) (string, error)
}

var modelStore map[string]*model

func init() {
	modelStore = make(map[string]*model)
}

func (m *model) addToRegistry(name string) {
	modelStore[name] = m
}

func (m *model) validate(s interface{}) (string, error) {
	val := reflect.ValueOf(s).Elem()
	for i, c := range m.constraints {
		v := c.validate(val.Field(i))
		if v != "" {
			return v, nil
		}
	}
	for _, v := range m.custom {
		v, err := v(s)
		if err != nil {
			return "", err
		}
		if v != "" {
			return v, nil
		}
	}
	return "", nil
}

func (m *model) registerStringConstraint(field reflect.StructField) {
	c := new(stringConstraint)
	c.field = strings.ToLower(field.Name)
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	c.max = m.getIntFromTag(field, tag, "max")
	c.min = m.getIntFromTag(field, tag, "min")
	in := m.getStringFromTag(tag, "in")
	if in != "" {
		c.in = strings.Split(in, ",")
	}
	reStr := m.getStringFromTag(tag, "regex")
	re, err := regexp.Compile(reStr)
	if err != nil {
		panic(fmt.Sprintf("govalid model registration error (%s.%s `regex:%s`): %s", m.name, field.Name, reStr, err.Error()))
	}
	c.regex = re
	m.constraints = append(m.constraints, c)
}

func (m *model) registerIntConstraint(field reflect.StructField) {
	c := new(intConstraint)
	c.field = strings.ToLower(field.Name)
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	c.max = m.getIntFromTag(field, tag, "max")
	c.min = m.getIntFromTag(field, tag, "min")
	in := m.getStringFromTag(tag, "in")
	if in != "" {
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
	c.field = strings.ToLower(field.Name)
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	c.max = m.getInt64FromTag(field, tag, "max")
	c.min = m.getInt64FromTag(field, tag, "min")
	in := m.getStringFromTag(tag, "in")
	if in != "" {
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
	c.field = strings.ToLower(field.Name)
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	c.max = m.getFloat64FromTag(field, tag, "max")
	c.min = m.getFloat64FromTag(field, tag, "min")
	m.constraints = append(m.constraints, c)
}

func (m *model) registerFloat32Constraint(field reflect.StructField) {
	c := new(float32Constraint)
	c.field = strings.ToLower(field.Name)
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = m.getBoolFromTag(tag, "req")
	c.max = m.getFloat32FromTag(field, tag, "max")
	c.min = m.getFloat32FromTag(field, tag, "min")
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

func (m *model) getIntFromTag(field reflect.StructField, tag string, key string) int {
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
			return i
		}
	}
	return 0
}

func (m *model) getInt64FromTag(field reflect.StructField, tag string, key string) int64 {
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
			return i
		}
	}
	return 0
}

func (m *model) getFloat32FromTag(field reflect.StructField, tag string, key string) float32 {
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
			return float32(f)
		}
	}
	return 0
}

func (m *model) getFloat64FromTag(field reflect.StructField, tag string, key string) float64 {
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
			return f
		}
	}
	return 0
}

func (m *model) getStringFromTag(tag string, key string) string {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			return parts[1]
		}
	}
	return ""
}
