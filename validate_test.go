package govalid_test

import (
	"testing"

	"github.com/twharmon/govalid"
)

func TestValidate(t *testing.T) {
	type nonPtrStruct struct {
		S string
	}
	govalid.Register(nonPtrStruct{})
	_, nonPtrErr := govalid.Validate(nonPtrStruct{})
	assertErr(t, "validate non-pointer", nonPtrErr)

	type Map map[string]string
	m := make(Map)
	_, mapPtrErr := govalid.Validate(&m)
	assertErr(t, "validate pointer to map", mapPtrErr)

	_, noRegisterErr := govalid.Validate(&struct{ S string }{"asdf"})
	assertErr(t, "validate unregistered pointer to struct", noRegisterErr)
}
