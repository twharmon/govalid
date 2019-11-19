package govalid_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/twharmon/govalid"
)

type constraintTest struct {
	Name     string
	Password string
}

func TestConstraint(t *testing.T) {
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
