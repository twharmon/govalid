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

	assertValid(t, "no validation rules with empty field", &n{})
	assertValid(t, "no validation rules with non-empty field", &n{b{"asdf"}})
}
