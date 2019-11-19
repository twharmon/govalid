package govalid

import (
	"reflect"
)

type constraint interface {
	violation(reflect.Value) error
	violations(reflect.Value) []error
}
