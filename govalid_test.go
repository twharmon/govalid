package govalid_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/twharmon/govalid"
)

type unsupported map[string]string

type registerTest struct {
	A unsupported
	S string
}

type constraintTest struct {
	Name     string
	Password string
}

func TestAddCustom(t *testing.T) {
	govalid.Register(constraintTest{})
	govalid.AddCustom(constraintTest{}, func(i interface{}) error {
		u := i.(*constraintTest)
		if u.Name == "" || u.Password == "" {
			return nil
		}
		if strings.Contains(u.Password, u.Name) {
			return errors.New("password can not contain name")
		}
		return nil
	})

	assertNoViolation(t, "custom validation rule with empty fields", &constraintTest{})
	assertViolation(t, "custom validation rule with invalid fields", &constraintTest{
		Name:     "Gopher",
		Password: "Gopher123",
	})

	assertPanic(t, "custom validation without registration", func() {
		govalid.AddCustom(struct{ S string }{}, func(i interface{}) error {
			return nil
		})
	})

	assertNoPanic(t, "custom validation using ptr", func() {
		type t struct{ S string }
		govalid.Register(t{})
		govalid.AddCustom(&t{}, func(i interface{}) error {
			return nil
		})
	})

	govalid.AddCustom(constraintTest{}, func(i interface{}) error {
		u := i.(*constraintTest)
		if u.Name == "error" {
			return errors.New("something went wrong during validation")
		}
		return nil
	})
	assertViolation(t, "custom validation rule with invalid fields", &constraintTest{
		Name: "error",
	})
}

func TestRegister(t *testing.T) {
	assertPanic(t, "when already registered", func() {
		govalid.Register(registerTest{})
		govalid.Register(registerTest{})
	})

	assertPanic(t, "when map", func() {
		type registerTestMap map[string]string
		testMap := make(registerTestMap)
		govalid.Register(testMap)
	})

	assertPanic(t, "when regex can't compile", func() {
		type testBadRegex struct {
			S string `validate:"regex:["`
		}
		govalid.Register(testBadRegex{})
	})

	assertPanic(t, "when int in won't parse", func() {
		type testBadIntIn struct {
			I int `validate:"in:4.3,5,6"`
		}
		govalid.Register(testBadIntIn{})
	})

	assertPanic(t, "when int64 in won't parse", func() {
		type testBadIntIn struct {
			I int `validate:"in:4.3,5,6"`
		}
		govalid.Register(testBadIntIn{})
	})

	assertPanic(t, "when int64 in won't parse", func() {
		type testBadInt64In struct {
			I int64 `validate:"in:4.3,5,6"`
		}
		govalid.Register(testBadInt64In{})
	})

	assertPanic(t, "when max int won't parse", func() {
		type testBadMaxInt struct {
			S string `validate:"max:4.3"`
		}
		govalid.Register(testBadMaxInt{})
	})

	assertPanic(t, "when max int64 won't parse", func() {
		type testBadMaxInt64 struct {
			I64 int64 `validate:"max:4.3"`
		}
		govalid.Register(testBadMaxInt64{})
	})

	assertPanic(t, "when max float32 won't parse", func() {
		type testBadMaxFloat32 struct {
			F32 float32 `validate:"max:4.g3"`
		}
		govalid.Register(testBadMaxFloat32{})
	})

	assertPanic(t, "when max float64 won't parse", func() {
		type testBadMaxFloat64 struct {
			F64 float64 `validate:"max:4.g3"`
		}
		govalid.Register(testBadMaxFloat64{})
	})

	assertNoPanic(t, "when ptr is registered", func() {
		type testBadMaxFloat64 struct {
			F64 float64 `validate:"max:4.3"`
		}
		govalid.Register(&testBadMaxFloat64{})
	})
}

func TestViolationNotStruct(t *testing.T) {
	var ty map[string]interface{}
	if govalid.Violation(ty) != govalid.ErrNotStruct {
		t.Fail()
	}
}

func TestViolationNotRegistered(t *testing.T) {
	var ty struct{}
	if govalid.Violation(ty) != govalid.ErrNotRegistered {
		t.Fail()
	}
}

func TestViolationsNotStruct(t *testing.T) {
	var ty map[string]interface{}
	if govalid.Violations(ty)[0] != govalid.ErrNotStruct {
		t.Fail()
	}
}

func TestViolationsNotRegistered(t *testing.T) {
	var ty struct{}
	if govalid.Violations(ty)[0] != govalid.ErrNotRegistered {
		t.Fail()
	}
}
