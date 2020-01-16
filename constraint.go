package govalid

import (
	"reflect"
)

type constraint interface {
	error(reflect.Value) error
	errors(reflect.Value) []error
}
