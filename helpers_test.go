package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func assertInvalid(t *testing.T, desc string, s interface{}) {
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error(err)
	}
	if v == "" {
		t.Errorf("assert invalid: %s (no violation)", desc)
	}
}

func assertValid(t *testing.T, desc string, s interface{}) {
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Errorf("assert valid: %s (violation: %s)", desc, v)
	}
}

func assertPanic(t *testing.T, desc string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("assert panic: %s (no panic)", desc)
		}
	}()
	f()
}

func assertErr(t *testing.T, desc string, err error) {
	if err == nil {
		t.Errorf("assert error: %s (nil error)", desc)
	}
}
