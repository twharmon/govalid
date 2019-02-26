package govalid_test

import (
	"strings"
	"testing"

	"github.com/twharmon/govalid"
)

type constraintTest struct {
	Name     string
	Password string
}

const errMsg = "password can not contain name"

func TestConstraint(t *testing.T) {
	govalid.Register(constraintTest{})
	govalid.AddCustom(constraintTest{}, func(i interface{}) (string, error) {
		u := i.(*constraintTest)
		if u.Name == "" || u.Password == "" {
			return "", nil
		}
		if strings.Contains(u.Password, u.Name) {
			return errMsg, nil
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
}
