# Govalid
Use Govalid to validate structs.
At load time, Govalid caches information about the structs you pass to `govalid.Register`.
Govalid uses that information at run time to reduce allocations and computing.

## Usage
Govalid currently supports the following constraints.

`int`
- req:"[bool]" - `req:"true"`
- min:"[int]" - `min:"5"`
- max:"[int]" - `max:"50"`
- in:"[int],[int],..." - `in:"1,2,3,4,5"`

`int64`
- req:"[bool]" - `req:"true"`
- min:"[int64]" - `min:"5"`
- max:"[int64]" - `max:"50"`
- in:"[int64],[int64],..." - `in:"1,2,3,4,5"`


`string`
- req:"[bool]" - `req:"true"`
- min:"[int]" - `min:"5"`
- max:"[int]" - `max:"50"`
- in:"[string],[string],..." - `in:"user,editor,admin"`
- regex:"[string]" - `regex:"^[a-zA-Z0-9]+$"`

See [examples](https://github.com/twharmon/govalid/tree/master/examples)
