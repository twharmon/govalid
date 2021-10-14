package govalid_test

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/twharmon/govalid"
)

func TestValidatorNotStruct(t *testing.T) {
	v := govalid.New()
	m := make(map[string]string)
	_, err := v.Violation(&m)
	equals(t, err, govalid.ErrNotStruct)
	_, err = v.Violations(&m)
	equals(t, err, govalid.ErrNotStruct)
}

func TestValidatorCustomValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar string `govalid:"req"`
	}
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
	foo := Foo{}
	vio, err := v.Violation(&foo)
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringCustomRuleValid(t *testing.T) {
	v := govalid.New()
	v.AddCustomStringRule("email", func(field string, value string) string {
		if strings.Contains(value, "@") {
			return ""
		}
		return fmt.Sprintf("%s must contain @", field)
	})
	type Foo struct {
		Bar string `govalid:"req|email"`
	}
	vio, err := v.Violation(&Foo{Bar: "baz@example.com"})
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&Foo{Bar: "baz@example.com"})
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorStringCustomRuleInvalid(t *testing.T) {
	v := govalid.New()
	v.AddCustomStringRule("email", func(field string, value string) string {
		if strings.Contains(value, "@") {
			return ""
		}
		return fmt.Sprintf("%s must contain @", field)
	})
	type Foo struct {
		Bar string `govalid:"req|email"`
	}
	vio, err := v.Violation(&Foo{Bar: "baz"})
	check(t, err)
	equals(t, vio, "Bar must contain @")
	vios, err := v.Violations(&Foo{Bar: "baz"})
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorInt64CustomRuleValid(t *testing.T) {
	v := govalid.New()
	v.AddCustomInt64Rule("even", func(field string, value int64) string {
		if value%2 == 0 {
			return ""
		}
		return fmt.Sprintf("%s must be even", field)
	})
	type Foo struct {
		Bar int64 `govalid:"req|even"`
	}
	vio, err := v.Violation(&Foo{Bar: 4})
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&Foo{Bar: 4})
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorInt64CustomRuleInvalid(t *testing.T) {
	v := govalid.New()
	v.AddCustomInt64Rule("even", func(field string, value int64) string {
		if value%2 == 0 {
			return ""
		}
		return fmt.Sprintf("%s must be even", field)
	})
	type Foo struct {
		Bar int64 `govalid:"req|even"`
	}
	vio, err := v.Violation(&Foo{Bar: 5})
	check(t, err)
	equals(t, vio, "Bar must be even")
	vios, err := v.Violations(&Foo{Bar: 5})
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorFloat64CustomRuleValid(t *testing.T) {
	v := govalid.New()
	v.AddCustomFloat64Rule("near-zero", func(field string, value float64) string {
		if math.Abs(value) < 1 {
			return ""
		}
		return fmt.Sprintf("%s must be near zero", field)
	})
	type Foo struct {
		Bar float64 `govalid:"req|near-zero"`
	}
	vio, err := v.Violation(&Foo{Bar: 0.1})
	check(t, err)
	equals(t, vio, "")
	vios, err := v.Violations(&Foo{Bar: 0.1})
	check(t, err)
	equals(t, len(vios), 0)
}

func TestValidatorFloat64CustomRuleInvalid(t *testing.T) {
	v := govalid.New()
	v.AddCustomFloat64Rule("near-zero", func(field string, value float64) string {
		t.Log(value, math.Abs(value))
		if math.Abs(value) < 1 {
			return ""
		}
		return fmt.Sprintf("%s must be near zero", field)
	})
	type Foo struct {
		Bar float64 `govalid:"req|near-zero"`
	}
	vio, err := v.Violation(&Foo{Bar: 5})
	check(t, err)
	equals(t, vio, "Bar must be near zero")
	vios, err := v.Violations(&Foo{Bar: 5})
	check(t, err)
	equals(t, len(vios), 1)
}

func TestValidatorIntReqValid(t *testing.T) {
	v := govalid.New()
	type Foo struct {
		Bar int `govalid:"req"`
	}
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
	foo := Foo{sql.NullString{Valid: true, String: "baz"}}
	vio, err := v.Violation(&foo)
	check(t, err)
	contains(t, vio, "Bar", "must be in")
	vios, err := v.Violations(&foo)
	check(t, err)
	equals(t, len(vios), 1)
}

func BenchmarkValidatorStringReqInvalid(b *testing.B) {
	v := govalid.New()
	type User struct {
		Name string `govalid:"req"`
	}
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
