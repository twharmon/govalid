package govalid

import (
	"log"
	"testing"
)

type user struct {
	ID             int64
	Name           string `req:"true" min:"5" max:"15" regex:"^[a-zA-Z]+$"`
	Email          string `req:"true" min:"3" max:"25" regex:".+?@.+?"`
	Age            int    `req:"false" min:"3" max:"120"`
	Role           string `req:"true" in:"admin,user"`
	FavoriteNumber int64  `req:"true" min:"1" max:"999999999999999"`
}

func init() {
	if err := Register(&user{}); err != nil {
		log.Fatalln(err)
	}
}

// BenchmarkGovalid .
func BenchmarkGovalid(b *testing.B) {
	u := &user{
		ID:             5,
		Name:           "Gopher",
		Email:          "gopher@example.com",
		Age:            45,
		Role:           "user",
		FavoriteNumber: 918273645,
	}
	for n := 0; n < b.N; n++ {
		Validate(u)
	}
}
