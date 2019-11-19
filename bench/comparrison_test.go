package bench

import (
	"testing"

	"github.com/twharmon/govalid"
)

type user struct {
	ID             int64
	Name           string  `validate:"req|min:5|max:15|regex:^[a-zA-Z]+$"`
	Email          string  `validate:"req|min:3|max:25|regex:^.+@.+$"`
	Age            int     `validate:"min:3|max:120"`
	Role           string  `validate:"in:admin,user,editor"`
	FavoriteNumber int64   `validate:"req|min:1|max:999999999999999"`
	Score          float64 `validate:"req|min:3.33|max:10.45"`
}

func init() {
	govalid.Register(user{})
}

func BenchmarkViolationsPass(b *testing.B) {
	user := &user{
		ID:             5,
		Name:           "Gopher",
		Email:          "admin@gmail.com",
		Age:            20,
		Role:           "admin",
		FavoriteNumber: 918273645,
		Score:          5.325,
	}
	for n := 0; n < b.N; n++ {
		govalid.Violations(user)
	}
}

func BenchmarkViolationFail(b *testing.B) {
	user := &user{
		ID:             5,
		Name:           "Goph",
		Email:          "admin@gmail.com",
		Age:            20,
		Role:           "admin",
		FavoriteNumber: 918273645,
		Score:          5.325,
	}
	for n := 0; n < b.N; n++ {
		govalid.Violation(user)
	}
}

func BenchmarkViolationsFail(b *testing.B) {
	user := &user{
		ID:             5,
		Name:           "Goph",
		Email:          "admingmail.com",
		Age:            2,
		Role:           "super_admin",
		FavoriteNumber: 918273645,
		Score:          50.325,
	}
	for n := 0; n < b.N; n++ {
		govalid.Violations(user)
	}
}
