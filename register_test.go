package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

type registerTest struct {
	S string
}

func TestRegister(t *testing.T) {
	assertPanic(t, "when register pointer", func() {
		govalid.Register(&registerTest{})
	})
}
