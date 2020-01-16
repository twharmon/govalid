package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func TestRegister(t *testing.T) {
	v := govalid.New()
	type T struct{}
	equals(t, v.Register(T{}), nil)
}

func TestRegisterDuplicate(t *testing.T) {
	v := govalid.New()
	type T struct{}
	check(t, v.Register(T{}))
	notEqual(t, v.Register(T{}), nil)
}
