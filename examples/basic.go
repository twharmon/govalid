package main

import (
	"fmt"
	"strings"

	"github.com/twharmon/govalid"
)

// User has information about a user
type User struct {
	ID             int64
	Name           string  `validate:"req|min:5|max:15|regex:^[a-zA-Z]+$"`
	Email          string  `validate:"req|min:3|max:25|regex:^.+?@.+?$"`
	Age            int     `validate:"req|min:3|max:120"`
	Role           string  `validate:"req|in:admin,user,editor"`
	FavoriteNumber int64   `validate:"req|min:1|max:999999999999999"`
	Score          float32 `validate:"req|min:3.33|max:10.45"`
	PreciseScore   float64 `validate:"req|min:1.5|max:6.22"`
}

// Post has information about a post
type Post struct {
	ID    int64
	Title string `validate:"req|min:1|max:30"`
	Body  string `validate:"req|min:100|max:65535"`
}

func main() {
	// Register the structs at load time that you will be validating later
	govalid.Register(User{}, Post{})

	// Add a custom validator
	govalid.AddCustom(User{}, func(obj interface{}) (string, error) {
		user := obj.(*User)
		if !strings.HasPrefix(user.Email, "admin@") && user.Role == "admin" {
			return "admin's email must start with 'admin@'", nil
		}
		return "", nil
	})

	user := &User{
		ID:             5,
		Name:           "Gopher",
		Email:          "user@gmail.com",
		Age:            2,
		Role:           "admin",
		FavoriteNumber: 918273645,
		Score:          4.325,
		PreciseScore:   3.325,
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
