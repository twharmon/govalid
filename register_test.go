package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

type unsupported map[string]string

type registerTest struct {
	A unsupported
	S string
}

func TestRegister(t *testing.T) {
	assertPanic(t, "when pointer", func() {
		govalid.Register(&registerTest{})
	})

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
}
