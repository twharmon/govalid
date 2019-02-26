package govalid_test

import (
	"testing"
	"time"

	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(n{})
}

type n struct {
	T time.Time
}

func TestNil(t *testing.T) {
	assertValid(t, "no validation rules with empty field", &n{})
	assertValid(t, "no validation rules with non-empty field", &n{time.Now()})
}
