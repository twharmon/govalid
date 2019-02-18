package govalid

import (
	"strings"
	"testing"
)

type user struct {
	ID             int64
	Name           string  `govalid:"req,min:5,max:15,regex:^[a-zA-Z]+$"`
	Email          string  `govalid:"req,min:3,max:25,regex:.+@.+"`
	Age            int     `govalid:"min:3,max:120"`
	Role           string  `govalid:"in:admin,user"`
	FavoriteNumber int64   `govalid:"req,min:1,max:999999999999999"`
	Score          float32 `govalid:"req,min:3.33,max:10.45"`
	PreciseScore   float64 `govalid:"req,min:1.5,max:6.22"`
}

func init() {
	Register(user{})
	AddCustom(user{}, func(obj interface{}) []string {
		user := obj.(*user)
		var violations []string
		if !strings.HasPrefix(user.Email, "admin@") && user.Role == "admin" {
			violations = append(violations, "admin's email must start with 'admin@'")
		}
		return violations
	})
}

// BenchmarkGovalid .
func BenchmarkGovalid(b *testing.B) {
	u := &user{
		ID:             5,
		Name:           "Gopher",
		Email:          "admin@gmail.com",
		Age:            2,
		Role:           "admin",
		FavoriteNumber: 918273645,
		Score:          5.325,
		PreciseScore:   5.325,
	}
	for n := 0; n < b.N; n++ {
		Validate(u)
	}
}
