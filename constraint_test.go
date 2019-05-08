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
	govalid.AddCustom(constraintTest{}, func(i interface{}) (string, error) {
		u := i.(*constraintTest)
		if u.Name == "" || u.Password == "" {
			return "", nil
		}
		if strings.Contains(u.Password, u.Name) {
			return "password can not contain name", nil
		}
		return "", nil
	})

	assertValid(t, "custom validation rule with empty fields", &constraintTest{})
	assertInvalid(t, "custom validation rule with invalid fields", &constraintTest{
		Name:     "Gopher",
		Password: "Gopher123",
	})

	assertPanic(t, "custom validation with pointer", func() {
		govalid.AddCustom(&constraintTest{}, func(i interface{}) (string, error) {
			return "", nil
		})
	})

	assertPanic(t, "custom validation without registration", func() {
		govalid.AddCustom(struct{ S string }{}, func(i interface{}) (string, error) {
			return "", nil
		})
	})

	govalid.AddCustom(constraintTest{}, func(i interface{}) (string, error) {
		u := i.(*constraintTest)
		if u.Name == "error" {
			return "", errors.New("something went wrong during validation")
		}
		return "", nil
	})
	_, err := govalid.Validate(&constraintTest{
		Name: "error",
	})
	assertErr(t, "custom validation error", err)
}
