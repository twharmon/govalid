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
	I int `govalid:"req"`
}

type iMin struct {
	I int `govalid:"min:5"`
}

type iReqMax struct {
	I int `govalid:"req|max:5"`
}

type iIn struct {
	I int `govalid:"in:1,2,3"`
}

func TestInt(t *testing.T) {
	assertNoViolation(t, "no validation rules with empty field", &i{})
	assertNoViolation(t, "no validation rules with non-empty field", &i{5})

	assertViolation(t, "`req` with empty field", &iReq{})
	assertNoViolation(t, "`req` with non-empty field", &iReq{5})

	assertNoViolation(t, "`min` with empty field", &iMin{})
	assertViolation(t, "`min` with field too less", &iMin{3})
	assertNoViolation(t, "`min` with valid field", &iMin{5})

	assertViolation(t, "`req|max` with empty field", &iReqMax{})
	assertViolation(t, "`req|max` with field too great", &iReqMax{7})
	assertNoViolation(t, "`req|max` with valid field", &iReqMax{3})

	assertNoViolation(t, "`in` with empty field", &iIn{})
	assertViolation(t, "`in` with invalid field", &iIn{7})
	assertNoViolation(t, "`in` with valid field", &iIn{3})
}
