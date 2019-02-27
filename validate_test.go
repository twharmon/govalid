package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func init() {
	govalid.Register(validateTestStruct{})
}

type validateTestStruct struct {
	S string
}

type validateTestMap map[string]string

func TestValidate(t *testing.T) {
	_, nonPtrErr := govalid.Validate(validateTestStruct{})
	assertErr(t, "validate non-pointer", nonPtrErr)

	testMap := make(validateTestMap)
	testMap["asdf"] = "fdsa"
	_, mapPtrErr := govalid.Validate(&testMap)
	assertErr(t, "validate pointer to map", mapPtrErr)

	_, noRegisterErr := govalid.Validate(&struct{ S string }{"asdf"})
	assertErr(t, "validate unregistered pointer to struct", noRegisterErr)
}
