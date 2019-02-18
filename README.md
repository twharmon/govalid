# Govalid
Use Govalid to validate structs.
At load time, Govalid caches information about the structs you pass to `govalid.Register`.
Govalid uses that information at run time to reduce allocations and computing.

## Usage
Govalid currently supports the following constraints.

`int`
- req:"[bool]" - `govalid:"req"`
- min:"[int]" - `govalid:"min:5"`
- max:"[int]" - `govalid:"max:50"`
- in:"[int],[int],..." - `govalid:"in:3,4,5"`

`int64`
- req:"[bool]" - `govalid:"req"`
- min:"[int64]" - `govalid:"min:5"`
- max:"[int64]" - `govalid:"max:50"`
- in:"[int64],[int64],..." - `govalid:"in:3,4,5"`

`float32`
- req:"[bool]" - `govalid:"req"`
- min:"[float32]" - `govalid:"min:5.0"`
- max:"[float32]" - `govalid:"max:50.0"`

`float64`
- req:"[bool]" - `govalid:"req"`
- min:"[float64]" - `govalid:"min:5.0"`
- max:"[float64]" - `govalid:"max:50.0"`

`string`
- req:"[bool]" - `govalid:"req"`
- min:"[int]" - `govalid:"min:5"`
- max:"[int]" - `govalid:"max:50"`
- in:"[string],[string],..." - `govalid:"in:user,editor,admin"`
- regex:"[string]" - `govalid:"regex:^[a-zA-Z0-9]+$"`

## Examples
See [examples](https://github.com/twharmon/govalid/tree/master/examples).

## Contribute
Create a pull request to contribute to Govalid.