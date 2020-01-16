package govalid_test

import (
	"errors"
	"testing"

	"github.com/twharmon/govalid"
)

func TestValidatorErrorNotRegistered(t *testing.T) {
	v := govalid.New()
	type T struct{}
	err := v.Error(&T{})
	equals(t, err, govalid.ErrNotRegistered)
}

func TestValidatorErrorNotPtr(t *testing.T) {
	v := govalid.New()
	type T struct{}
	check(t, v.Register(T{}))
	err := v.Error(T{})
	equals(t, err, govalid.ErrNotPtrToStruct)
}

func TestValidatorErrorStringReqInvalid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"req"`
	}
	check(t, v.Register(T{}))
	err := v.Error(&T{})
	equals(t, errors.Is(err, govalid.ErrInvalidStruct), true)
}
