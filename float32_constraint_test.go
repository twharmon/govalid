package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(f32{}, f32Req{}, f32ReqMin{}, f32Max{})
}

type f32 struct {
	F32 float32
}

type f32Req struct {
	F32 float32 `validate:"req"`
}

type f32ReqMin struct {
	F32 float32 `validate:"req|min:5.5"`
}

type f32Max struct {
	F32 float32 `validate:"max:5.5"`
}

func TestFloat32(t *testing.T) {
	assertValid(t, "no validation rules with empty field", &f32{})
	assertValid(t, "no validation rules with non-empty field", &f32{5})

	assertInvalid(t, "`req` with empty field", &f32Req{})
	assertValid(t, "`req` with non-empty field", &f32Req{5.5})

	assertInvalid(t, "`req|min` with empty field", &f32ReqMin{})
	assertInvalid(t, "`req|min` with field too less", &f32ReqMin{3.5})
	assertValid(t, "`req|min` with valid field", &f32ReqMin{5.5})

	assertValid(t, "`max` with empty field", &f32Max{})
	assertInvalid(t, "`max` with field too great", &f32Max{7.5})
	assertValid(t, "`max` with valid field", &f32Max{3.5})
}
