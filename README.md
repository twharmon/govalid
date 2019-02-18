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
Here is an example struct with some basic constraints.
```
type User struct {
	ID             int64
	Name           string  `govalid:"req,min:5,max:25,regex:^[a-zA-Z ]+$"`
	Email          string  `govalid:"req,min:3,max:100,regex:.+?@.+?"`
	Age            int     `govalid:"min:18,max:120"`
	Role           string  `govalid:"in:admin,editor,user"`
	Grade          float32 `govalid:"min:0.0,max:100.0"`
}
```

Register the struct and add a custom constraint.
Both Register and AddCustom can panic, so they should be called at load time.
```
govalid.Register(User{})
govalid.AddCustom(User{}, func(obj interface{}) []string {
    user := obj.(*User)
    var violations []string
    if !strings.HasPrefix(user.Email, "admin@") && user.Role == "admin" {
        violations = append(violations, "admin's email must start with 'admin@'")
    }
    return violations
})
```

Validate a struct.
```
user := &User{
    ID:    5,
    Name:  "Gopher",
    Email: "user@gmail.com",
    Age:   0,
    Role:  "admin",
    Grade: 87.5,
}
userViolations, err := govalid.Validate(user)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(userViolations) // ["age must be at least 18", "admin's email must start with 'admin'"]
```

See [examples](https://github.com/twharmon/govalid/tree/master/examples) for more.

## Contribute
Create a pull request to contribute to Govalid.