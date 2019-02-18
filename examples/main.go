package main

import (
	"fmt"
	"log"

	"github.com/twharmon/govalid"
)

// User has information about a user
type User struct {
	ID             int64
	Name           string `req:"true" min:"5" max:"15" regex:"^[a-zA-Z]+$"`
	Email          string `req:"true" min:"3" max:"25" regex:".+?@.+?"`
	Age            int    `req:"false" min:"3" max:"120"`
	Role           string `req:"true" in:"admin,user"`
	FavoriteNumber int64  `req:"true" min:"1" max:"999999999999999"`
}

// Post has information about a post
type Post struct {
	ID    int64
	Title string `req:"true" min:"1" max:"30"`
	Body  string `req:"true" min:"100" max:"65535"`
}

func main() {
	// Register the structs at load time that you will be validating later.
	if registerErr := govalid.Register(&User{}, &Post{}); registerErr != nil {
		log.Fatalln(registerErr)
	}

	user := &User{
		ID:             5,
		Name:           "Gopher",
		Email:          "gopher@example.com",
		Age:            45,
		Role:           "user",
		FavoriteNumber: 918273645,
	}
	userErr := govalid.Validate(user)
	fmt.Println(userErr) // <nil>

	post := &Post{
		ID:    1,
		Title: "Hello, World!",
		Body:  "Hello!",
	}
	postErr := govalid.Validate(post)
	fmt.Println(postErr) // Body must be at least 100 characters
}
