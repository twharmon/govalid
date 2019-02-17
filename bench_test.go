package govalid

import (
	"log"
	"testing"
)

type user struct {
	ID     int64  `json:"id"`
	Name   string `json:"name" req:"no" min:"5" max:"15"`
	Email  string `json:"email" req:"true" min:"3" max:"25" regex:"@"`
	Age    int    `json:"age" req:"true" min:"18" max:"120"`
	Active bool   `json:"active"`
	Test   int64  `json:"test" req:"true" min:"1" max:"8793465928238947520"`
}

func init() {
	if err := MustPrepare(&user{}); err != nil {
		log.Fatalln(err)
	}
}

// BenchmarkValidate .
func BenchmarkValidate(b *testing.B) {
	u := new(user)
	u.ID = 5
	u.Name = "Gopher"
	u.Email = "gopher@example.com"
	u.Age = 44
	u.Active = true
	u.Test = 44
	for n := 0; n < b.N; n++ {
		MustValidate(u)
	}
}
