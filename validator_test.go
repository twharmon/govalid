package govalid_test

import (
	"fmt"
	"testing"

	"github.com/twharmon/govalid"
)

func TestValidatorViolationNotRegistered(t *testing.T) {
	v := govalid.New()
	type T struct{}
	_, err := v.Violation(&T{})
	equals(t, err, govalid.ErrNotRegistered)
}

func TestValidatorViolationNotPtr(t *testing.T) {
	v := govalid.New()
	type T struct{}
	_, err := v.Violation(T{})
	equals(t, err, govalid.ErrNotPtrToStruct)
}

func TestValidatorViolationNotStruct(t *testing.T) {
	v := govalid.New()
	m := make(map[string]string)
	_, err := v.Violation(&m)
	equals(t, err, govalid.ErrNotPtrToStruct)
}

func TestValidatorViolationCustomValid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"req"`
	}
	check(t, v.Register(T{}))
	check(t, v.AddCustom(T{}, func(v interface{}) string {
		if v.(*T).F == "bar" {
			return "F must not be bar"
		}
		return ""
	}))
	vio, err := v.Violation(&T{"foo"})
	check(t, err)
	equals(t, vio, "")
}

func TestValidatorViolationCustomInvalid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"req"`
	}
	check(t, v.Register(T{}))
	check(t, v.AddCustom(T{}, func(v interface{}) string {
		if v.(*T).F == "foo" {
			return "F must not be foo"
		}
		return ""
	}))
	vio, err := v.Violation(&T{"foo"})
	check(t, err)
	equals(t, vio, "F must not be foo")
}

func TestValidatorViolationStringReqInvalid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"req"`
	}
	check(t, v.Register(T{}))
	vio, err := v.Violation(&T{})
	check(t, err)
	contains(t, vio, "F", "required")
}

func TestValidatorViolationStringReqValid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"req"`
	}
	check(t, v.Register(T{}))
	vio, err := v.Violation(&T{"foo"})
	check(t, err)
	equals(t, vio, "")
}

func TestValidatorViolationStringRegexInvalid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"regex:^[a-z]+$"`
	}
	check(t, v.Register(T{}))
	vio, err := v.Violation(&T{"Foo"})
	check(t, err)
	contains(t, vio, "F", "match", "regex")
}

func TestValidatorViolationStringRegexValid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"regex:^[a-z]+$"`
	}
	check(t, v.Register(T{}))
	vio, err := v.Violation(&T{"foo"})
	check(t, err)
	equals(t, vio, "")
}

func TestValidatorViolationStringMinInvalid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"min:6"`
	}
	check(t, v.Register(T{}))
	vio, err := v.Violation(&T{"foo"})
	check(t, err)
	contains(t, vio, "F", "at least")
}

func TestValidatorViolationStringMinValid(t *testing.T) {
	v := govalid.New()
	type T struct {
		F string `govalid:"min:6"`
	}
	check(t, v.Register(T{}))
	vio, err := v.Violation(&T{"foobarbaz"})
	check(t, err)
	equals(t, vio, "")
}

func ExampleValidator() {
	v := govalid.New()
	type User struct {
		Name string `govalid:"req|max:8"`
	}
	v.Register(User{})
	vio, _ := v.Violation(&User{"foobarbaz"})
	fmt.Println(vio)
	// Output: Name can not be longer than 8 characters
}

func BenchmarkValidatorViolationStringReqInvalid(b *testing.B) {
	v := govalid.New()
	type User struct {
		Name string `govalid:"req"`
	}
	v.Register(User{})
	user := User{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Violation(&user)
	}
}

func BenchmarkValidatorViolationStringReqValid(b *testing.B) {
	v := govalid.New()
	type User struct {
		Name string `govalid:"req"`
	}
	v.Register(User{})
	user := User{"foo"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Violation(&user)
	}
}

func BenchmarkValidatorViolationsInValid(b *testing.B) {
	v := govalid.New()
	type User struct {
		Name string `govalid:"req|min:2|max:32|regex:[a-z]+"`
		Role string `govalid:"req|in:user,editor,admin"`
		Age  int    `govalid:"req|min:18"`
	}
	v.Register(User{})
	user := User{
		Name: "foo",
		Role: "super_admin",
		Age:  10,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Violations(&user)
	}
}
