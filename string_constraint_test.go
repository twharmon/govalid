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
	assertValid(t, "no validation rules with empty field", &str{})
	assertValid(t, "no validation rules with non-empty field", &str{"asdf", "asdf"})

	assertInvalid(t, "`req` with empty field", &strReq{})
	assertValid(t, "`req` with non-empty field", &strReq{"asdf"})

	assertValid(t, "`min` with empty field", &strMin{})
	assertInvalid(t, "`min` with field too short", &strMin{"asdf"})
	assertValid(t, "`min` with valid field", &strMin{"asdfasdf"})

	assertInvalid(t, "`req|max` with empty field", &strReqMax{})
	assertInvalid(t, "`req|max` with field too long", &strReqMax{"fdsa", "asdfasdf"})
	assertValid(t, "`req|max` with valid field", &strReqMax{"fdsa", "asdf"})

	assertValid(t, "`regex` with empty field", &strRegex{})
	assertInvalid(t, "`regex` with invalid field", &strRegex{"asdf0"})
	assertValid(t, "`regex` with valid field", &strRegex{"asdf"})

	assertValid(t, "`in` with empty field", &strIn{})
	assertInvalid(t, "`in` with invalid field", &strIn{"abcd"})
	assertValid(t, "`in` with valid field", &strIn{"def"})

	assertValid(t, "`max` with empty struct field", &strNullMax{})
	assertInvalid(t, "`max` with struct field too long", &strNullMax{sql.NullString{Valid: true, String: "asdfasdf"}})
	assertValid(t, "`max` with valid struct field", &strNullMax{sql.NullString{Valid: true, String: "asdf"}})
	assertValid(t, "`max` with valid quirky struct field", &strNullMax{sql.NullString{Valid: false, String: "asdasdfasdff"}})
}
