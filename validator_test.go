package govalid_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/twharmon/govalid"
)

func TestValidatorNotRegistered(t *testing.T) {
	v := govalid.New()
	type Foo struct{}
	_, err := v.Violation(&Foo{})
	equals(t, err, govalid.ErrNotRegistered)
	_, err = v.Violations(&Foo{})
	equals(t, err, govalid.ErrNotRegistered)
}

func TestValidatorCustomNotRegistered(t *testing.T) {
	v := govalid.New()
	type Foo struct{}
	err := v.AddCustom(Foo{}, func(v interface{}) string {
		return ""
	})
	notEqual(t, err, nil)
}

func TestValidatorNotPtr(t *testing.T) {
	v := govalid.New()
	type Foo struct{}
	_, err := v.Violation(Foo{})
	equals(t, err, govalid.ErrNotPtrToStruct)
	_, err = v.Violations(Foo{})
	equals(t, err, govalid.ErrNotPtrToStruct)
}

func TestValidatorNotStruct(t *testing.T) {
	v := govalid.New()
	m := make(map[string]string)
	_, err := v.Violation(&m)
	equals(t, err, govalid.ErrNotPtrToStruct)
	_, err = v.Violations(&m)
	equals(t, err, govalid.ErrNotPtrToStruct)
}

func TestValidatorRegisterNotStruct(t *testing.T) {
	v := govalid.New()
	m := make(map[string]string)
	notEqual(t, v.Register(m), nil)
}

func TestValidatorRegisterPointer(t *testing.T) {
	v := govalid.New()
	type T struct{}
	equals(t, v.Register(&T{}), nil)
}

func TestValidatorRegisterNilConstraint(t *testing.T) {
	v := govalid.New()
	type T struct {
		A struct{}
		b string
		M map[string]string
	}
	equals(t, v.Register(T{}), nil)
}

func TestValidatorCustomValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	check(t, v.AddCustom(Foo{}, func(v interface{}) string {
		if v.(*Foo).Bar == "bar" {
			return "Bar must not be bar"
		}
		return ""
	}))
	foo := Foo{"foo"}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorCustomPtrValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	check(t, v.AddCustom(&Foo{}, func(v interface{}) string {
		if v.(*Foo).Bar == "bar" {
			return "Bar must not be bar"
		}
		return ""
	}))
	foo := Foo{"foo"}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorCustomInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	check(t, v.AddCustom(Foo{}, func(v interface{}) string {
		if v.(*Foo).Bar == "foo" {
			return "Bar must not be foo"
		}
		return ""
	}))
	foo := Foo{"foo"}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "Bar must not be foo")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorNilConstraint(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		S   string
		I   int
		I64 int64
		F32 float32
		F64 float64
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorIntReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{1}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorIntReqInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "required")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorIntMinValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorIntMinInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{5}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "at least")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorIntMinNotReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorIntMaxValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req|max:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorIntMaxInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req|max:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{11}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "can not be greater than")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorIntInValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req|in:10,20"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorIntInInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req|in:10,20"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{11}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "must be in")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorInt64ReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{1}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorInt64ReqInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "required")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorInt64MinValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorInt64MinInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{5}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "at least")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorInt64MinNotReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorInt64MaxValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"req|max:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorInt64MaxInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"req|max:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{11}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "can not be greater than")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorInt64InValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"req|in:10,20"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorInt64InInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"req|in:10,20"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{11}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "must be in")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}
func TestValidatorNullInt64Invalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar sql.NullInt64 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{sql.NullInt64{Valid: true, Int64: 0}}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "must be at least")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorNullInt64Valid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar sql.NullInt64 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{sql.NullInt64{Valid: true, Int64: 10}}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat32ReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{1}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat32ReqInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "required")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorFloat32MinValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat32MinInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{5}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "at least")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorFloat32MinNotReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat32MaxValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"req|max:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat32MaxInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"req|max:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{11}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "can not be greater than")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorFloat64ReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{1}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat64ReqInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "required")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorFloat64MinValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat64MinInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{5}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "at least")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorFloat64MinNotReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat64MaxValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"req|max:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{10}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat64MaxInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"req|max:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{11}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "can not be greater than")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorNullFloat64Invalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar sql.NullFloat64 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{sql.NullFloat64{Valid: true, Float64: 0}}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "must be at least")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorNullFloat64Valid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar sql.NullFloat64 `govalid:"req|min:10"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{sql.NullFloat64{Valid: true, Float64: 10}}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringReqInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "required")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorStringReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"req"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"foo"}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringRegexInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"regex:^[a-z]+$"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"Foo"}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "match", "regex")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorStringRegexValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"regex:^[a-z]+$"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"foo"}
	vio, err := v.Violation(&Foo{"foo"})
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringMinInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"req|min:6"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"foo"}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "at least")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorStringMinValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"req|min:6"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"foobarbaz"}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringMinNotReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"min:6"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{""}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringMaxInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"max:6"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"foobarbaz"}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "longer")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorStringMaxValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"max:6"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"foo"}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringInValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"in:foo,bar"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"foo"}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringInInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"in:foo,bar"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{"baz"}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "must be in")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorNullStringInvalid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar sql.NullString `govalid:"in:foo,bar"`
	}
	check(t, v.Register(Foo{}))
	foo := Foo{sql.NullString{Valid: true, String: "baz"}}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "must be in")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorStringInvalidMaxTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"max:10.5"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorNullStringInvalidMaxTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar sql.NullString `govalid:"max:10.5"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorFloat32InvalidMaxTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"max:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorFloat64InvalidMaxTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"max:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorNullFloat64InvalidMaxTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar sql.NullFloat64 `govalid:"max:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorIntInvalidMaxTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"max:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorInt64InvalidMaxTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"max:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorNullInt64InvalidMaxTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar sql.NullInt64 `govalid:"max:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorStringInvalidMinTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"min:10.5"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorInt64InvalidRegexTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"regex:^([)][[[[[]$"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorFloat32InvalidMinTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float32 `govalid:"min:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorFloat64InvalidMinTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar float64 `govalid:"min:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorIntInvalidMinTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"min:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorInt64InvalidMinTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"min:foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorIntInvalidInTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"in:5,foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
}

func TestValidatorInt64InvalidInTag(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int64 `govalid:"in:5,foo"`
	}
	notEqual(t, v.Register(Foo{}), nil)
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

func BenchmarkValidatorStringReqInvalid(b *testing.B) {
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

func BenchmarkValidatorStringReqValid(b *testing.B) {
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

func BenchmarkValidatorsVariety(b *testing.B) {
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
