package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

type unsupported map[string]string

type registerTest struct {
	A unsupported
	S string
}

func TestRegister(t *testing.T) {
	assertPanic(t, "when pointer", func() {
		govalid.Register(&registerTest{})
	})

	assertPanic(t, "when already registered", func() {
		govalid.Register(registerTest{})
		govalid.Register(registerTest{})
	})

	assertPanic(t, "when map", func() {
		type registerTestMap map[string]string
		testMap := make(registerTestMap)
		govalid.Register(testMap)
	})
}
