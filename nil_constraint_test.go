package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func TestNil(t *testing.T) {
	type b struct {
		B string
	}
	type n struct {
		B b
	}
	govalid.Register(n{})

	assertNilViolation(t, "no validation rules with empty field", &n{})
	assertNilViolation(t, "no validation rules with non-empty field", &n{b{"asdf"}})
}
