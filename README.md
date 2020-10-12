# Govalid

![](https://github.com/twharmon/govalid/workflows/Test/badge.svg) [![](https://goreportcard.com/badge/github.com/twharmon/govalid)](https://goreportcard.com/report/github.com/twharmon/govalid) [![](https://gocover.io/_badge/github.com/twharmon/govalid)](https://gocover.io/github.com/twharmon/govalid) [![](https://pkg.go.dev/github.com/twharmon/govalid)](https://pkg.go.dev/github.com/twharmon/govalid)

Use Govalid to validate structs.

## Documentation

For full documentation see [pkg.go.dev](https://pkg.go.dev/github.com/twharmon/govalid).

## Example

```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/twharmon/govalid"
)

type Post struct {
	// ID has no constraints
	ID int

	// Title is required, must be at least 3 characters long, cannot be
	// more than 20 characters long, and must match ^[a-zA-Z ]+$
	Title string `govalid:"req|min:3|max:20|regex:^[a-zA-Z ]+$"`

	// Body is not required, and cannot be more than 10000 charachers.
	Body string `govalid:"max:10000"`

	// Category is not required, but if not zero value ("") it must be
	// either "announcement" or "bookreview".
	Category string `govalid:"in:announcement,bookreview"`
}

var v = govalid.New()

func main() {
	v.Register(Post{}) // Register all structs at load time

	// Add Custom validation to `Post`
	v.AddCustom(Post{}, func(val interface{}) string {
		post := val.(*Post)
		if post.Category != "" && !strings.Contains(post.Body, post.Category) {
			return fmt.Sprintf("Body must contain %s", post.Category)
		}
		return ""
	})

	p := Post{
		ID:       5,
		Title:    "Hi",
		Body:     "World",
		Category: "announcement",
	}

	vio, err := v.Violations(&p)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(vio)
}
```

## Benchmarks

```
BenchmarkValidatorStringReqInvalid	         192 ns/op	      48 B/op	       3 allocs/op
BenchmarkValidatorStringReqValid	        98.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkValidatorsVariety	                1114 ns/op	     281 B/op	      13 allocs/op
```

## Contribute

Make a pull request.
