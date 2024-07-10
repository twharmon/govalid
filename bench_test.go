package govalid_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/twharmon/govalid"
)

func BenchmarkValidateStringReqInvalid(b *testing.B) {
	type User struct {
		Name string `valid:"req"`
	}
	user := User{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		govalid.Validate(&user)
	}
}

func BenchmarkValidateStringReqValid(b *testing.B) {
	type User struct {
		Name string `valid:"req"`
	}
	user := User{"foo"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		govalid.Validate(&user)
	}
}

func BenchmarkValidateVariety(b *testing.B) {
	govalid.Rule("role", func(v any) (string, error) {
		switch tv := v.(type) {
		case string:
			if tv == "user" || tv == "editor" || tv == "admin" {
				return "", nil
			}
			return "role must be user, editor, or admin", nil
		default:
			return "", errors.New("role must be applied to string only")
		}
	})
	nameRegExp := regexp.MustCompile("[a-z]+")
	govalid.Rule("name", func(v any) (string, error) {
		switch tv := v.(type) {
		case string:
			if nameRegExp.MatchString(tv) {
				return "", nil
			}
			return "name alphanumeric", nil
		default:
			return "", errors.New("name must be applied to string only")
		}
	})
	type User struct {
		Name string `valid:"req|min:2|max:32|name"`
		Role string `valid:"req|role"`
		Age  int    `valid:"req|min:18"`
	}
	user := User{
		Name: "foo",
		Role: "super_admin",
		Age:  10,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		govalid.Validate(&user)
	}
}
