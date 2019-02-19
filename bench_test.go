package govalid

import (
	"regexp"
	"strings"
	"testing"

	validator "gopkg.in/go-playground/validator.v9"
)

type user struct {
	ID             int64
	Name           string   `validate:"req|min:5|max:15|regex:^[a-zA-Z]+$"`
	Email          string   `validate:"req|min:3|max:25|regex:.+@.+"`
	Age            int      `validate:"min:3|max:120"`
	Role           string   `validate:"in:admin,user,editor"`
	FavoriteNumber int64    `validate:"req|min:1|max:999999999999999"`
	Score          float32  `validate:"req|min:3.33|max:10.45"`
	PreciseScore   float64  `validate:"req|min:1.5|max:6.22"`
	Tags           []string `validate:"req|min:1|max:5"`
}

type user2 struct {
	ID             int64
	Name           string  `validate:"required,gte=5,lte=15"`
	Email          string  `validate:"required,gte=3,lte=25,alpha"`
	Age            int     `validate:"gte=3,lte=120,custom-email"`
	Role           string  `validate:"role"`
	FavoriteNumber int64   `validate:"required,gte=1,lte=999999999999999"`
	Score          float32 `validate:"required,gte=3.33,lte=10.45"`
	PreciseScore   float64 `validate:"required,gte=1.5,lte=6.22"`
}

var validate *validator.Validate

func init() {
	Register(user{})
	AddCustom(user{}, func(obj interface{}) ([]string, error) {
		user := obj.(*user)
		var violations []string
		if !strings.HasPrefix(user.Email, "admin@") && user.Role == "admin" {
			violations = append(violations, "admin's email must start with 'admin@'")
		}
		return violations, nil
	})
	validate = validator.New()
	validate.RegisterValidation("alpha", func(fl validator.FieldLevel) bool {
		re, err := regexp.Compile("^[a-zA-Z]+$")
		if err != nil {
			panic(err)
		}
		return re.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return s == "admin" || s == "user" || s == "editor"
	})
	validate.RegisterValidation("custom-email", func(fl validator.FieldLevel) bool {
		re, err := regexp.Compile("^.+@.+$")
		if err != nil {
			panic(err)
		}
		return re.MatchString(fl.Field().String())
	})
}

// BenchmarkGovalid .
func BenchmarkGovalid(b *testing.B) {
	user := &user{
		ID:             5,
		Name:           "Goph er",
		Email:          "admin@gmail.com",
		Age:            2,
		Role:           "admin",
		FavoriteNumber: 918273645,
		Score:          5.325,
		PreciseScore:   5.325,
		Tags:           []string{"a", "b"},
	}
	for n := 0; n < b.N; n++ {
		Validate(user)
	}
}

// BenchmarkValidatorV9 .
func BenchmarkValidatorV9(b *testing.B) {
	user := &user2{
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
		validate.Struct(user)
	}
}
