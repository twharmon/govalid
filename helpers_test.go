package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func assertNoViolation(t *testing.T, desc string, s interface{}) {
	err := govalid.Violation(s)
	if err != nil {
		t.Errorf("assert nil violation: %s (found %s)", desc, err)
	}

	errs := govalid.Violations(s)
	if len(errs) > 0 {
		t.Errorf("assert nil violation: %s (found %s)", desc, errs)
	}
}

func assertViolation(t *testing.T, desc string, s interface{}) {
	err := govalid.Violation(s)
	if err == nil {
		t.Errorf("assert violation: %s (found none)", desc)
	}

	errs := govalid.Violations(s)
	if len(errs) == 0 {
		t.Errorf("assert violation: %s (found none)", desc)
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
