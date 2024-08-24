# Govalid

[![](https://goreportcard.com/badge/github.com/twharmon/govalid)](https://goreportcard.com/report/github.com/twharmon/govalid)

Use Govalid to validate structs.

## Documentation

For full documentation see [pkg.go.dev](https://pkg.go.dev/github.com/twharmon/govalid).

## Todo
Document use of `dive`. Rules are now ordered.
- pointers: dereference and check remaining rules
- slices / arrays: check each item with remaining rules
- structs: validates the struct according to its own field tags (remaining rules have no meaning)

## Example

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/twharmon/govalid"
)

type Post struct {
	// ID has no constraints
	ID int

	// Title is required, must be at least 3 characters long, and
	// cannot be more than 20 characters long
	Title string `valid:"req|min:3|max:20"`

	// Body is not required, cannot be more than 10000 charachers,
	// and must be "fun" (a custom rule defined below).
	Body string `valid:"max:10000|fun"`
}

func main() {
	// Add custom string rule "fun" that can be used on any string field
	// in any struct.
	govalid.Rule("fun", func(v any) (string, error) {
		switch tv := v.(type) {
		case string:
			if float64(strings.Count(tv, "!"))/float64(utf8.RuneCountInString(tv)) > 0.001 {
				return "", nil
			}
			return "must contain more exclamation marks", nil
		default:
			return "", errors.New("fun constraint must be applied to string only")
		}
	})

	p := Post{
		ID:    5,
		Title: "Hi",
		Body:  "Hello world!",
	}

	vio, err := govalid.Validate(&p)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(vio)
}
```

## Contribute

Make a pull request.
