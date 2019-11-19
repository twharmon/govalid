package main

import (
	"fmt"

	"github.com/twharmon/govalid"
)

// User contains user information
type User struct {
	ID    int64
	Name  string  `govalid:"req|min:5|max:25,regex:^[a-zA-Z ]+$"`
	Email string  `govalid:"req|min:3|max:100|regex:^.+?@.+?$"`
	Age   int     `govalid:"min:8|max:100"`
	Role  string  `govalid:"req|in:admin,editor,user"`
	Grade float32 `govalid:"min:0.0|max:100.0"`
}

func main() {
	govalid.Register(&User{})

	u := &User{
		ID:    5,
		Name:  "Gopher",
		Email: "gopher@example.com",
		Age:   11,
		Role:  "mascot",
		Grade: 99.5,
	}

	violation := govalid.Violation(u)
	fmt.Println("violation:", violation)
}
