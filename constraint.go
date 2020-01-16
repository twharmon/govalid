package govalid

import (
	"reflect"
)

type constraint interface {
	violation(reflect.Value) string
	violations(reflect.Value) []string
}
