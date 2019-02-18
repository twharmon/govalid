package main

import (
	"fmt"
	"strings"

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
	// Register the structs at load time that you will be validating later
	govalid.Register(User{}, Post{})

	// Add a custom validator
	govalid.AddCustom(User{}, func(obj interface{}) []string {
		user := obj.(*User)
		var violations []string
		if !strings.HasPrefix(user.Email, "admin") && user.Role == "admin" {
			violations = append(violations, "admin's email must start with 'admin'")
		}
		return violations
	})

	user := &User{
		ID:             5,
		Name:           "Gopher",
		Email:          "asdf@gmail.com",
		Age:            2,
		Role:           "admin",
		FavoriteNumber: 918273645,
	}
	userViolations, err := govalid.Validate(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(userViolations) // ["age must be at least 3", "admin's email must start with 'admin'"]

	post := &Post{
		ID:    1,
		Title: "Hello, World!",
		Body:  "Hello!",
	}
	postViolations, err := govalid.Validate(post)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(postViolations) // ["body must be at least 100 characters"]
}
