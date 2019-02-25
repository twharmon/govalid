package govalid_test

import (
	"strings"
	"testing"

	"github.com/twharmon/govalid"
)

type stringReqStruct struct {
	Name string `validate:"req|min:5|max:10"`
}

type stringStruct struct {
	Name string `validate:"min:5|max:10|regex:^[a-zA-Z0-9]+$"`
}

type stringInStruct struct {
	Name string `validate:"in:abc,def,ghi"`
}

func init() {
	govalid.Register(stringStruct{}, stringReqStruct{}, stringInStruct{})
}

func TestValidateStringReqMinInvalid(t *testing.T) {
	s := &stringReqStruct{"asdf"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if !strings.Contains(v, "at least") {
		t.Errorf("violation was '%s'; should contain 'at least'", v)
	}
}

func TestValidateStringReqMinInvalidOmit(t *testing.T) {
	s := &stringReqStruct{}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if !strings.Contains(v, "required") {
		t.Errorf("violation was '%s'; should contain 'required'", v)
	}
}

func TestValidateStringReqMinValid(t *testing.T) {
	s := &stringReqStruct{"12345"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if v != "" {
		t.Errorf("violation was '%s'; should be ''", v)
	}
}

func TestValidateStringMinInvalid(t *testing.T) {
	s := &stringStruct{"1234"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if !strings.Contains(v, "at least") {
		t.Errorf("violation was '%s'; should contain 'at least'", v)
	}
}

func TestValidateStringMinValid(t *testing.T) {
	s := &stringStruct{"12345"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if v != "" {
		t.Errorf("violation was '%s'; should be ''", v)
	}
}

func TestValidateStringMinValidOmit(t *testing.T) {
	s := &stringStruct{}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if v != "" {
		t.Errorf("violation was '%s'; should be ''", v)
	}
}

func TestValidateStringMaxInvalid(t *testing.T) {
	s := &stringStruct{"12345678901"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if !strings.Contains(v, "longer than") {
		t.Errorf("violation was '%s'; should contiain 'longer than'", v)
	}
}

func TestValidateStringReqMaxInvalid(t *testing.T) {
	s := &stringReqStruct{"12345678901"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if !strings.Contains(v, "longer than") {
		t.Errorf("violation was '%s'; should contain 'longer than'", v)
	}
}

func TestValidateStringRegexInvalid(t *testing.T) {
	s := &stringStruct{"1236 78901"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if !strings.Contains(v, "regex") {
		t.Errorf("violation was '%s'; should contain 'regex'", v)
	}
}

func TestValidateStringRegexValid(t *testing.T) {
	s := &stringStruct{"aV2asdf"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if v != "" {
		t.Errorf("violation was '%s'; should be ''", v)
	}
}

func TestValidateStringInInvalid(t *testing.T) {
	s := &stringInStruct{"asdf"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if !strings.Contains(v, "must be in") {
		t.Errorf("violation was '%s'; should contain 'must be in'", v)
	}
}

func TestValidateStringInValid(t *testing.T) {
	s := &stringInStruct{"def"}
	v, err := govalid.Validate(s)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if v != "" {
		t.Errorf("violation was '%s'; should be ''", v)
	}
}
