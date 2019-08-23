package govalid_test

import (
	"testing"
	"time"

	"github.com/twharmon/gosql"
	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(tm{}, tmMin{}, tmNullMax{})
}

type tm struct {
	T time.Time
}

type tmMin struct {
	T time.Time `validate:"min:0"`
}

type tmNullMax struct {
	T gosql.NullTime `validate:"req|max:3600"`
}

func TestTime(t *testing.T) {
	now := time.Now()
	assertNilViolation(t, "no validation rules with empty field", &tm{})
	assertNilViolation(t, "no validation rules with non-empty field", &tm{now})

	assertNilViolation(t, "`min` with empty field", &tmMin{})
	assertNilViolation(t, "`min` with valid field", &tmMin{now})
	assertViolation(t, "`min` with invalid field", &tmMin{now.Add(time.Hour * 10)})

	assertViolation(t, "`req|max` with empty struct field", &tmNullMax{})
	assertViolation(t, "`req|max` with invalid struct field", &tmNullMax{gosql.NullTime{Valid: true, Time: now.AddDate(0, 0, -1)}})
}
