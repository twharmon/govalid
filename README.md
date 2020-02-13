# Govalid

![](https://github.com/twharmon/govalid/workflows/Test/badge.svg) [![](https://goreportcard.com/badge/github.com/twharmon/govalid)](https://goreportcard.com/report/github.com/twharmon/govalid) [![](https://gocover.io/_badge/github.com/twharmon/govalid)](https://gocover.io/github.com/twharmon/govalid) [![GoDoc](https://godoc.org/github.com/twharmon/govalid?status.svg)](https://godoc.org/github.com/twharmon/govalid)

Use Govalid to validate structs.

## Documentation
For full documentation see [godoc](https://godoc.org/github.com/twharmon/govalid).

## Example
```
package main

import (
	"fmt"

	"github.com/twharmon/govalid"
)

type Post struct {
	ID    int
	Title string `govalid:"req|min:3|max:20"`
	Body  string `govalid:"max:10000"`
}

func main() {
	v := govalid.New()
	v.Register(Post{})
	p := Post{
		ID:    5,
		Title: "Hi",
		Body:  "World",
	}
	vio, _ := v.Violation(&p)
	fmt.Println(vio)
}
```

## Benchmarks
```
BenchmarkValidatorViolationStringReqInvalid-4   	 5856351	       202 ns/op	      48 B/op	       3 allocs/op
BenchmarkValidatorViolationStringReqValid-4     	13598485	        88.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkValidatorViolationsVariety-4           	 1000000	      1022 ns/op	     289 B/op	      14 allocs/op
```

## Contribute
Make a pull request.
