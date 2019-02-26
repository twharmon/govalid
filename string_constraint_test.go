package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(str{}, strReq{}, strMin{}, strReqMax{}, strRegex{}, strIn{})
}

type str struct {
	S string
}

type strReq struct {
	S string `validate:"req"`
}

type strMin struct {
	S string `validate:"min:5"`
}

type strReqMax struct {
	S string `validate:"req|max:5"`
}

type strRegex struct {
	S string `validate:"regex:^[a-z]+$"`
}

type strIn struct {
	S string `validate:"in:abc,def,ghi"`
}

func TestString(t *testing.T) {
	assertValid(t, "no validation rules with empty field", &str{})
	assertValid(t, "no validation rules with non-empty field", &str{"asdf"})

	assertInvalid(t, "`req` with empty field", &strReq{})
	assertValid(t, "`req` with non-empty field", &strReq{"asdf"})

	assertValid(t, "`min` with empty field", &strMin{})
	assertInvalid(t, "`min` with field too short", &strMin{"asdf"})
	assertValid(t, "`min` with valid field", &strMin{"asdfasdf"})

	assertInvalid(t, "`req|max` with empty field", &strReqMax{})
	assertInvalid(t, "`req|max` with field too long", &strReqMax{"asdfasdf"})
	assertValid(t, "`req|max` with valid field", &strReqMax{"asdf"})

	assertValid(t, "`regex` with empty field", &strRegex{})
	assertInvalid(t, "`regex` with invalid field", &strRegex{"asdf0"})
	assertValid(t, "`regex` with valid field", &strRegex{"asdf"})

	assertValid(t, "`in` with empty field", &strIn{})
	assertInvalid(t, "`in` with invalid field", &strIn{"abcd"})
	assertValid(t, "`in` with valid field", &strIn{"def"})
}
