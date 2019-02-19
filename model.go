package govalid

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagKey = "validate"

type model struct {
	constraints []constraint
	custom      []func(interface{}) ([]string, error)
}

var modelStore map[string]*model

func init() {
	modelStore = make(map[string]*model)
}

func (m *model) addToRegistry(name string) {
	modelStore[name] = m
}

func (m *model) validate(s interface{}) ([]string, error) {
	val := reflect.ValueOf(s).Elem()
	var vs []string
	for i, c := range m.constraints {
		vs = append(vs, c.validate(val.Field(i))...)
	}
	for _, v := range m.custom {
		newVs, err := v(s)
		if err != nil {
			return nil, err
		}
		vs = append(vs, newVs...)
	}
	return vs, nil
}

func (m *model) registerStringConstraint(field reflect.StructField) {
	c := new(stringConstraint)
	c.field = strings.ToLower(field.Name)
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return
	}
	c.req = getBoolFromTag(tag, "req")
	c.max = getIntFromTag(tag, "max")
	c.min = getIntFromTag(tag, "min")
	in := getStringFromTag(tag, "in")
	if in != "" {
		c.in = strings.Split(in, ",")
	}
	re, err := regexp.Compile(getStringFromTag(tag, "regex"))
	if err != nil {
		panic(err)
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
	c.req = getBoolFromTag(tag, "req")
	c.max = getIntFromTag(tag, "max")
	c.min = getIntFromTag(tag, "min")
	in := getStringFromTag(tag, "in")
	if in != "" {
		for _, optStr := range strings.Split(in, ",") {
			opt, err := strconv.Atoi(optStr)
			if err != nil {
				panic(err)
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
	c.req = getBoolFromTag(tag, "req")
	c.max = getInt64FromTag(tag, "max")
	c.min = getInt64FromTag(tag, "min")
	in := getStringFromTag(tag, "in")
	if in != "" {
		for _, optStr := range strings.Split(in, ",") {
			opt, err := strconv.ParseInt(optStr, 10, 64)
			if err != nil {
				panic(err)
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
	c.req = getBoolFromTag(tag, "req")
	c.max = getFloat64FromTag(tag, "max")
	c.min = getFloat64FromTag(tag, "min")
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
	c.req = getBoolFromTag(tag, "req")
	c.max = getFloat32FromTag(tag, "max")
	c.min = getFloat32FromTag(tag, "min")
	m.constraints = append(m.constraints, c)
}

func (m *model) registerNilConstraint() {
	c := new(nilConstraint)
	m.constraints = append(m.constraints, c)
}
