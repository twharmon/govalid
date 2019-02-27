package govalid_test

import (
	"testing"
	"time"

	"github.com/twharmon/govalid"
)

func TestNil(t *testing.T) {
	type n struct {
		T time.Time
	}
	govalid.Register(n{})

	assertValid(t, "no validation rules with empty field", &n{})
	assertValid(t, "no validation rules with non-empty field", &n{time.Now()})
}
