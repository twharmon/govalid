package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(i64{}, i64Req{}, i64Min{}, i64Max{}, i64ReqIn{})
}

type i64 struct {
	I64 int64
}

type i64Req struct {
	I64 int64 `validate:"req"`
}

type i64Min struct {
	I64 int64 `validate:"min:5"`
}

type i64Max struct {
	I64 int64 `validate:"max:5"`
}

type i64ReqIn struct {
	I64 int64 `validate:"req|in:1,2,3"`
}

func TestInt64(t *testing.T) {
	assertValid(t, "no validation rules with empty field", &i64{})
	assertValid(t, "no validation rules with non-empty field", &i64{5})

	assertInvalid(t, "`req` with empty field", &i64Req{})
	assertValid(t, "`req` with non-empty field", &i64Req{5})

	assertValid(t, "`min` with empty field", &i64Min{})
	assertInvalid(t, "`min` with field too less", &i64Min{3})
	assertValid(t, "`min` with valid field", &i64Min{5})

	assertValid(t, "`max` with empty field", &i64Max{})
	assertInvalid(t, "`max` with field too great", &i64Max{7})
	assertValid(t, "`max` with valid field", &i64Max{3})

	assertInvalid(t, "`req|in` with empty field", &i64ReqIn{})
	assertInvalid(t, "`req|in` with invalid field", &i64ReqIn{7})
	assertValid(t, "`req|in` with valid field", &i64ReqIn{3})
}
