package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(i{}, iReq{}, iMin{}, iReqMax{}, iIn{})
}

type i struct {
	I int
}

type iReq struct {
	I int `validate:"req"`
}

type iMin struct {
	I int `validate:"min:5"`
}

type iReqMax struct {
	I int `validate:"req|max:5"`
}

type iIn struct {
	I int `validate:"in:1,2,3"`
}

func TestInt(t *testing.T) {
	assertNilViolation(t, "no validation rules with empty field", &i{})
	assertNilViolation(t, "no validation rules with non-empty field", &i{5})

	assertViolation(t, "`req` with empty field", &iReq{})
	assertNilViolation(t, "`req` with non-empty field", &iReq{5})

	assertNilViolation(t, "`min` with empty field", &iMin{})
	assertViolation(t, "`min` with field too less", &iMin{3})
	assertNilViolation(t, "`min` with valid field", &iMin{5})

	assertViolation(t, "`req|max` with empty field", &iReqMax{})
	assertViolation(t, "`req|max` with field too great", &iReqMax{7})
	assertNilViolation(t, "`req|max` with valid field", &iReqMax{3})

	assertNilViolation(t, "`in` with empty field", &iIn{})
	assertViolation(t, "`in` with invalid field", &iIn{7})
	assertNilViolation(t, "`in` with valid field", &iIn{3})
}
