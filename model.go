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
	custom      []func(interface{}) string
}

func (m *model) violation(s interface{}) string {
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i, c := range m.constraints {
		if v := c.violation(val.Field(i)); v != "" {
			return v
		}
	}
	for _, v := range m.custom {
		if msg := v(s); msg != "" {
			return msg
		}
	}
	return ""
}

func (m *model) violations(s interface{}) []string {
	var vs []string
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i, c := range m.constraints {
		vs = append(vs, c.violations(val.Field(i))...)
	}
	for _, v := range m.custom {
		if msg := v(s); msg != "" {
			vs = append(vs, msg)
		}
	}
	return vs
}

func (m *model) registerStringConstraint(field reflect.StructField, customRules map[string]func(string, string) string) error {
	c := new(stringConstraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return nil
	}
	var allowedRuleNames = []string{"req", "min", "max", "in", "regex"}
	for rule := range customRules {
		allowedRuleNames = append(allowedRuleNames, rule)
	}
	if err := m.getInvalidTagErr(tag, allowedRuleNames...); err != nil {
		return err
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok, err := m.getIntFromTag(field, tag, "max"); err != nil {
		return err
	} else if ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok, err := m.getIntFromTag(field, tag, "min"); err != nil {
		return err
	} else if ok {
		c.isMinSet = true
		c.min = min
	}
	if in, ok := m.getStringFromTag(tag, "in"); ok {
		c.in = strings.Split(in, ",")
	}
	if reStr, ok := m.getStringFromTag(tag, "regex"); ok {
		re, err := regexp.Compile(reStr)
		if err != nil {
			return fmt.Errorf("govalid model registration error (%s.%s `regex:%s`): %w", m.name, field.Name, reStr, err)
		}
		c.regex = re
	}
	for ruleName := range customRules {
		if m.getBoolFromTag(tag, ruleName) {
			c.customRules = append(c.customRules, customRules[ruleName])
		}
	}
	m.constraints = append(m.constraints, c)
	return nil
}

func (m *model) registerIntConstraint(field reflect.StructField) error {
	c := new(intConstraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return nil
	}
	if err := m.getInvalidTagErr(tag, "req", "min", "max", "in"); err != nil {
		return err
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok, err := m.getIntFromTag(field, tag, "max"); err != nil {
		return err
	} else if ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok, err := m.getIntFromTag(field, tag, "min"); err != nil {
		return err
	} else if ok {
		c.isMinSet = true
		c.min = min
	}
	if in, ok := m.getStringFromTag(tag, "in"); ok {
		for _, optStr := range strings.Split(in, ",") {
			opt, err := strconv.Atoi(optStr)
			if err != nil {
				return fmt.Errorf("govalid model registration error (%s.%s `in:%s`): %w", m.name, field.Name, in, err)
			}
			c.in = append(c.in, opt)
		}
	}
	m.constraints = append(m.constraints, c)
	return nil
}

func (m *model) registerInt64Constraint(field reflect.StructField, customRules map[string]func(string, int64) string) error {
	c := new(int64Constraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return nil
	}
	var allowedRuleNames = []string{"req", "min", "max", "in"}
	for rule := range customRules {
		allowedRuleNames = append(allowedRuleNames, rule)
	}
	if err := m.getInvalidTagErr(tag, allowedRuleNames...); err != nil {
		return err
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok, err := m.getInt64FromTag(field, tag, "max"); err != nil {
		return err
	} else if ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok, err := m.getInt64FromTag(field, tag, "min"); err != nil {
		return err
	} else if ok {
		c.isMinSet = true
		c.min = min
	}
	if in, ok := m.getStringFromTag(tag, "in"); ok {
		for _, optStr := range strings.Split(in, ",") {
			opt, err := strconv.ParseInt(optStr, 10, 64)
			if err != nil {
				return fmt.Errorf("govalid model registration error (%s.%s `in:%s`): %w", m.name, field.Name, in, err)
			}
			c.in = append(c.in, opt)
		}
	}
	for ruleName := range customRules {
		if m.getBoolFromTag(tag, ruleName) {
			c.customRules = append(c.customRules, customRules[ruleName])
		}
	}
	m.constraints = append(m.constraints, c)
	return nil
}

func (m *model) registerFloat64Constraint(field reflect.StructField, customRules map[string]func(string, float64) string) error {
	c := new(float64Constraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return nil
	}
	var allowedRuleNames = []string{"req", "min", "max"}
	for rule := range customRules {
		allowedRuleNames = append(allowedRuleNames, rule)
	}
	if err := m.getInvalidTagErr(tag, allowedRuleNames...); err != nil {
		return err
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok, err := m.getFloat64FromTag(field, tag, "max"); err != nil {
		return err
	} else if ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok, err := m.getFloat64FromTag(field, tag, "min"); err != nil {
		return err
	} else if ok {
		c.isMinSet = true
		c.min = min
	}
	for ruleName := range customRules {
		if m.getBoolFromTag(tag, ruleName) {
			c.customRules = append(c.customRules, customRules[ruleName])
		}
	}
	m.constraints = append(m.constraints, c)
	return nil
}

func (m *model) registerFloat32Constraint(field reflect.StructField) error {
	c := new(float32Constraint)
	c.field = field.Name
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		m.registerNilConstraint()
		return nil
	}
	if err := m.getInvalidTagErr(tag, "req", "min", "max", "in"); err != nil {
		return err
	}
	c.req = m.getBoolFromTag(tag, "req")
	if max, ok, err := m.getFloat32FromTag(field, tag, "max"); err != nil {
		return err
	} else if ok {
		c.isMaxSet = true
		c.max = max
	}
	if min, ok, err := m.getFloat32FromTag(field, tag, "min"); err != nil {
		return err
	} else if ok {
		c.isMinSet = true
		c.min = min
	}
	m.constraints = append(m.constraints, c)
	return nil
}

func (m *model) registerNilConstraint() {
	c := new(nilConstraint)
	m.constraints = append(m.constraints, c)
}

func (m *model) getInvalidTagErr(tag string, allowedRuleNames ...string) error {
	for _, rule := range strings.Split(tag, "|") {
		ruleName := strings.Split(rule, ":")[0]
		found := false
		for _, allowedName := range allowedRuleNames {
			if ruleName == allowedName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid rule '%s' found in govalid tag on %s", ruleName, m.name)
		}
	}
	return nil
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

func (m *model) getIntFromTag(field reflect.StructField, tag string, key string) (int, bool, error) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			i, err := strconv.Atoi(parts[1])
			if err != nil {
				return 0, false, fmt.Errorf("govalid model registration error (%s.%s `%s`): %w", m.name, field.Name, c, err)
			}
			return i, true, nil
		}
	}
	return 0, false, nil
}

func (m *model) getInt64FromTag(field reflect.StructField, tag string, key string) (int64, bool, error) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			i, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return 0, false, fmt.Errorf("govalid model registration error (%s.%s `%s`): %w", m.name, field.Name, c, err)
			}
			return i, true, nil
		}
	}
	return 0, false, nil
}

func (m *model) getFloat32FromTag(field reflect.StructField, tag string, key string) (float32, bool, error) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			f, err := strconv.ParseFloat(parts[1], 32)
			if err != nil {
				return 0, false, fmt.Errorf("govalid model registration error (%s.%s `%s`): %w", m.name, field.Name, c, err)
			}
			return float32(f), true, nil
		}
	}
	return 0, false, nil
}

func (m *model) getFloat64FromTag(field reflect.StructField, tag string, key string) (float64, bool, error) {
	cs := strings.Split(tag, "|")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			f, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return 0, false, fmt.Errorf("govalid model registration error (%s.%s `%s`): %w", m.name, field.Name, c, err)
			}
			return f, true, nil
		}
	}
	return 0, false, nil
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
