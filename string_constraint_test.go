package govalid_test

import (
	"database/sql"
	"testing"

	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(str{}, strReq{}, strMin{}, strReqMax{}, strRegex{}, strIn{}, strNullMax{})
}

type str struct {
	ignoredField string
	S            string
}

type strReq struct {
	S string `validate:"req"`
}

type strMin struct {
	S string `validate:"min:5"`
}

type strReqMax struct {
	ignoredField string
	S            string `validate:"req|max:5"`
}

type strRegex struct {
	S string `validate:"regex:^[a-z]+$"`
}

type strIn struct {
	S string `validate:"in:abc,def,ghi"`
}

type strNullMax struct {
	S sql.NullString `validate:"max:5"`
}

func TestString(t *testing.T) {
	assertNilViolation(t, "no validation rules with empty field", &str{})
	assertNilViolation(t, "no validation rules with non-empty field", &str{"asdf", "asdf"})

	assertViolation(t, "`req` with empty field", &strReq{})
	assertNilViolation(t, "`req` with non-empty field", &strReq{"asdf"})

	assertNilViolation(t, "`min` with empty field", &strMin{})
	assertViolation(t, "`min` with field too short", &strMin{"asdf"})
	assertNilViolation(t, "`min` with valid field", &strMin{"asdfasdf"})

	assertViolation(t, "`req|max` with empty field", &strReqMax{})
	assertViolation(t, "`req|max` with field too long", &strReqMax{"fdsa", "asdfasdf"})
	assertNilViolation(t, "`req|max` with valid field", &strReqMax{"fdsa", "asdf"})

	assertNilViolation(t, "`regex` with empty field", &strRegex{})
	assertViolation(t, "`regex` with invalid field", &strRegex{"asdf0"})
	assertNilViolation(t, "`regex` with valid field", &strRegex{"asdf"})

	assertNilViolation(t, "`in` with empty field", &strIn{})
	assertViolation(t, "`in` with invalid field", &strIn{"abcd"})
	assertNilViolation(t, "`in` with valid field", &strIn{"def"})

	assertNilViolation(t, "`max` with empty struct field", &strNullMax{})
	assertViolation(t, "`max` with struct field too long", &strNullMax{sql.NullString{Valid: true, String: "asdfasdf"}})
	assertNilViolation(t, "`max` with valid struct field", &strNullMax{sql.NullString{Valid: true, String: "asdf"}})
	assertNilViolation(t, "`max` with valid quirky struct field", &strNullMax{sql.NullString{Valid: false, String: "asdasdfasdff"}})
}
