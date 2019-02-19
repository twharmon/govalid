# Govalid
Use Govalid to validate structs.
At load time, Govalid caches information about the structs you pass to `govalid.Register`.
Govalid uses that information at run time to reduce allocations and computing.

## Usage
Govalid currently supports the following constraints.

`int` & `int64`
- req
- min
- max
- in

`float32` & `float64`
- req
- min
- max

`string`
- req
- min
- max
- in
- regex

## Examples
Here is an example struct with some basic constraints.
```
type User struct {
    ID             int64
    Name           string  `validate:"req|min:5|max:25,regex:^[a-zA-Z ]+$"`
    Email          string  `validate:"req|min:3|max:100|regex:^.+?@.+?$"`
    Age            int     `validate:"min:18|max:120"`
    Role           string  `validate:"req|in:admin,editor,user"`
    Grade          float32 `validate:"min:0.0|max:100.0"`
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
    // two errors are possible
    // 1) you did not register User yet (govalid.Register(User{}))
    // 2) you did not pass a pointer to User
    fmt.Println(err)
    return
}
fmt.Println(userViolations) // ["age must be at least 18", "admin's email must start with 'admin'"]
```

See [examples](https://github.com/twharmon/govalid/tree/master/examples) for more.

## Contribute
Create a pull request to contribute to Govalid.