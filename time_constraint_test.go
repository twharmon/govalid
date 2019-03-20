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
	assertValid(t, "no validation rules with empty field", &tm{})
	assertValid(t, "no validation rules with non-empty field", &tm{now})

	assertValid(t, "`min` with empty field", &tmMin{})
	assertValid(t, "`min` with valid field", &tmMin{now})
	assertInvalid(t, "`min` with invalid field", &tmMin{now.Add(time.Hour * 10)})

	assertInvalid(t, "`req|max` with empty struct field", &tmNullMax{})
	assertInvalid(t, "`req|max` with invalid struct field", &tmNullMax{gosql.NullTime{Valid: true, Time: now.AddDate(0, 0, -1)}})
}
