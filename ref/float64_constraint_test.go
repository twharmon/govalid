package govalid_test

import (
	"database/sql"
	"testing"

	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(f64{}, f64Req{}, f64ReqMin{}, f64Max{}, f64NullMax{})
}

type f64 struct {
	F64 float64
}

type f64Req struct {
	F64 float64 `govalid:"req"`
}

type f64ReqMin struct {
	F64 float64 `govalid:"req|min:5.5"`
}

type f64Max struct {
	F64 float64 `govalid:"max:5.5"`
}

type f64NullMax struct {
	F64 sql.NullFloat64 `govalid:"max:5.5"`
}

func TestFloat64(t *testing.T) {
	assertNoViolation(t, "no validation rules with empty field", &f64{})
	assertNoViolation(t, "no validation rules with non-empty field", &f64{5})

	assertViolation(t, "`req` with empty field", &f64Req{})
	assertNoViolation(t, "`req` with non-empty field", &f64Req{5.5})

	assertViolation(t, "`req|min` with empty field", &f64ReqMin{})
	assertViolation(t, "`req|min` with field too less", &f64ReqMin{3.5})
	assertNoViolation(t, "`req|min` with valid field", &f64ReqMin{5.5})

	assertNoViolation(t, "`max` with empty field", &f64Max{})
	assertViolation(t, "`max` with field too great", &f64Max{7.5})
	assertNoViolation(t, "`max` with valid field", &f64Max{3.5})

	assertNoViolation(t, "`max` with empty struct field", &f64NullMax{})
	assertViolation(t, "`max` with struct field too great", &f64NullMax{sql.NullFloat64{Valid: true, Float64: 7.5}})
	assertNoViolation(t, "`max` with valid struct field", &f64NullMax{sql.NullFloat64{Valid: true, Float64: 3.5}})
	assertNoViolation(t, "`max` with valid quirky struct field", &f64NullMax{sql.NullFloat64{Valid: false, Float64: 13.5}})
}
