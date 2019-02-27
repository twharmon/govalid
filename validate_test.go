package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func TestValidate(t *testing.T) {
	type validateTestStruct struct {
		S string
	}
	govalid.Register(validateTestStruct{})
	_, nonPtrErr := govalid.Validate(validateTestStruct{})
	assertErr(t, "validate non-pointer", nonPtrErr)

	type validateTestMap map[string]string
	testMap := make(validateTestMap)
	_, mapPtrErr := govalid.Validate(&testMap)
	assertErr(t, "validate pointer to map", mapPtrErr)

	_, noRegisterErr := govalid.Validate(&struct{ S string }{"asdf"})
	assertErr(t, "validate unregistered pointer to struct", noRegisterErr)
}
