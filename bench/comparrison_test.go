package bench

import (
	"regexp"
	"testing"

	"github.com/twharmon/govalid"
	validator "gopkg.in/go-playground/validator.v9"
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

type user2 struct {
	ID             int64
	Name           string  `validate:"required,gte=5,lte=15,alpha"`
	Email          string  `validate:"required,gte=3,lte=25,custom-email"`
	Age            int     `validate:"gte=3,lte=120"`
	Role           string  `validate:"role"`
	FavoriteNumber int64   `validate:"required,gte=1,lte=999999999999999"`
	Score          float64 `validate:"required,gte=3.33,lte=10.45"`
}

var validate *validator.Validate

func init() {
	govalid.Register(user{})
	validate = validator.New()
	reAl, err := regexp.Compile("^[a-zA-Z]+$")
	if err != nil {
		panic(err)
	}
	validate.RegisterValidation("alpha", func(fl validator.FieldLevel) bool {
		return reAl.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return s == "admin" || s == "user" || s == "editor"
	})
	reEm, err := regexp.Compile("^.+@.+$")
	if err != nil {
		panic(err)
	}
	validate.RegisterValidation("custom-email", func(fl validator.FieldLevel) bool {
		return reEm.MatchString(fl.Field().String())
	})
}

func BenchmarkGovalidPass(b *testing.B) {
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

func BenchmarkValidatorV9Pass(b *testing.B) {
	user := &user2{
		ID:             5,
		Name:           "Gopher",
		Email:          "admin@gmail.com",
		Age:            20,
		Role:           "admin",
		FavoriteNumber: 918273645,
		Score:          5.325,
	}
	for n := 0; n < b.N; n++ {
		validate.Struct(user)
	}
}

func BenchmarkGovalidOneError(b *testing.B) {
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

func BenchmarkValidatorV9OneError(b *testing.B) {
	user := &user2{
		ID:             5,
		Name:           "Goph",
		Email:          "admin@gmail.com",
		Age:            20,
		Role:           "admin",
		FavoriteNumber: 918273645,
		Score:          5.325,
	}
	for n := 0; n < b.N; n++ {
		validate.Struct(user)
	}
}

func BenchmarkGovalidManyErrors(b *testing.B) {
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

func BenchmarkValidatorV9ManyErrors(b *testing.B) {
	user := &user2{
		ID:             5,
		Name:           "Goph",
		Email:          "admingmail.com",
		Age:            2,
		Role:           "super_admin",
		FavoriteNumber: 918273645,
		Score:          50.325,
	}
	for n := 0; n < b.N; n++ {
		validate.Struct(user)
	}
}
