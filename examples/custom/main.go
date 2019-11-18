package main

import (
	"errors"
	"log"

	"github.com/twharmon/govalid"
)

// User contains user information
type User struct {
	ID   int64
	Name string `validate:"req"`
}

func main() {
	govalid.Register(User{})
	govalid.AddCustom(User{}, func(obj interface{}) error {
		user := obj.(*User)
		if user.Name == "Gopher" {
			return errors.New("No gophers allowed")
		}
		return nil
	})

	u := &User{
		ID:   5,
		Name: "Gopher",
	}

	violations := govalid.Violations(u)
	if len(violations) > 0 {
		log.Fatalln(violations)
	}
}
