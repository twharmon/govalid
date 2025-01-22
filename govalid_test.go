package govalid_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/twharmon/govalid"
)

func TestValidateArgument(t *testing.T) {
	type A struct {
		A string `valid:"req"`
	}
	t.Run("fail: struct", func(t *testing.T) {
		errMustInclude(t, A{}, "required", "A")
	})
	t.Run("fail: pointer", func(t *testing.T) {
		errMustInclude(t, &A{}, "required", "A")
	})
	t.Run("fail: map", func(t *testing.T) {
		errMustInclude(t, map[string]any{}, "of kind")
	})
	t.Run("fail: nil", func(t *testing.T) {
		errMustInclude(t, nil, "nil")
	})
}

func TestValidateUnexported(t *testing.T) {
	type A struct {
		a string
	}
	t.Run("ok", func(t *testing.T) {
		errMustBeNil(t, A{a: ""})
	})
}

func TestValidateIllegalUse(t *testing.T) {
	type A struct {
		A string `valid:"req|min:foo"`
	}
	t.Run("min", func(t *testing.T) {
		errMustInclude(t, A{A: "a"}, "A")
	})
}

func TestValidateString(t *testing.T) {
	t.Run("fail: req", func(t *testing.T) {
		errMustInclude(t, struct {
			A string `valid:"req"`
		}{}, "required", "A")
	})
	t.Run("ok: req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A string `valid:"req"`
		}{A: "a"})
	})
	t.Run("fail: min", func(t *testing.T) {
		errMustInclude(t, struct {
			A string `valid:"min:3"`
		}{A: "ab"}, "min", "3")
	})
	t.Run("ok: min", func(t *testing.T) {
		errMustBeNil(t, struct {
			A string `valid:"min:3"`
		}{A: "abc"})
	})
	t.Run("ok: min not req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A string `valid:"min:3"`
		}{})
	})
	t.Run("fail: max", func(t *testing.T) {
		errMustInclude(t, struct {
			A string `valid:"max:3"`
		}{A: "abcd"}, "max", "3")
	})
	t.Run("ok: max", func(t *testing.T) {
		errMustBeNil(t, struct {
			A string `valid:"max:3"`
		}{A: "abc"})
	})
}

func TestValidateStringSlice(t *testing.T) {
	t.Run("fail: req", func(t *testing.T) {
		errMustInclude(t, struct {
			A []string `valid:"req"`
		}{}, "required", "A")
	})
	t.Run("ok: req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A []string `valid:"req"`
		}{A: []string{}})
	})
	t.Run("ok: req|dive|req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A []string `valid:"req|dive|req"`
		}{A: []string{"a"}})
	})
	t.Run("fail: req|dive|req", func(t *testing.T) {
		errMustInclude(t, struct {
			A []string `valid:"req|dive|req"`
		}{A: []string{"a", ""}}, "required", "A", "index", "1")
	})
	t.Run("fail: min:3", func(t *testing.T) {
		errMustInclude(t, struct {
			A []string `valid:"min:3"`
		}{A: []string{"a", "b"}}, "min", "3")
	})
	t.Run("ok: min:3", func(t *testing.T) {
		errMustBeNil(t, struct {
			A []string `valid:"min:3"`
		}{A: []string{"a", "b", "c"}})
	})
	t.Run("fail: max:3", func(t *testing.T) {
		errMustInclude(t, struct {
			A []string `valid:"max:3"`
		}{A: []string{"a", "b", "c", "d"}}, "max", "3")
	})
	t.Run("ok: max:3", func(t *testing.T) {
		errMustBeNil(t, struct {
			A []string `valid:"max:3"`
		}{A: []string{"a", "b", "c"}})
	})
	t.Run("ok: min:3 not req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A []string `valid:"min:3"`
		}{})
	})
	t.Run("fail: min:3 not req", func(t *testing.T) {
		errMustInclude(t, struct {
			A []string `valid:"min:3"`
		}{A: []string{}}, "min", "3")
	})
	t.Run("ok: more elems than rules (panic check)", func(t *testing.T) {
		errMustBeNil(t, struct {
			A []string `valid:"dive|min:1"`
		}{A: []string{"a", "b", "c"}})
	})
	t.Run("ok: more elems than rules (panic check)", func(t *testing.T) {
		errMustInclude(t, struct {
			A []string `valid:"dive|min:1|req"`
		}{A: []string{"a", "b", ""}}, "required")
	})
}

func TestValidatePointerToStringSlice(t *testing.T) {
	t.Run("fail: req", func(t *testing.T) {
		errMustInclude(t, struct {
			A []*string `valid:"req"`
		}{}, "required", "A")
	})
	t.Run("ok: req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A []*string `valid:"req"`
		}{A: []*string{}})
	})
	t.Run("ok: req|dive|req|dive|req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A []*string `valid:"req|dive|req|dive|req"`
		}{A: []*string{ptr("a")}})
	})
	t.Run("fail: req|dive|req|dive|req", func(t *testing.T) {
		errMustInclude(t, struct {
			A []*string `valid:"req|dive|req|dive|req"`
		}{A: []*string{ptr("a"), ptr("")}}, "required", "A", "index", "1")
	})
}

func TestValidateStructSlice(t *testing.T) {
	type A struct {
		A string `valid:"req"`
	}
	type B struct {
		As []A `valid:"req|dive"`
	}
	t.Run("fail: req", func(t *testing.T) {
		errMustInclude(t, B{}, "required", "As")
	})
	t.Run("ok: req", func(t *testing.T) {
		errMustBeNil(t, B{As: []A{}})
	})
	t.Run("fail: slice items", func(t *testing.T) {
		errMustInclude(t, B{As: []A{{}}}, "required", "index", "0", "A")
	})
}

func TestValidatePointerToString(t *testing.T) {
	t.Run("fail: req", func(t *testing.T) {
		errMustInclude(t, struct {
			A *string `valid:"req"`
		}{}, "required", "A")
	})
	t.Run("ok: req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A *string `valid:"req"`
		}{A: ptr("a")})
	})
	t.Run("ok: req|dive|req", func(t *testing.T) {
		errMustBeNil(t, struct {
			A *string `valid:"req|dive|req"`
		}{A: ptr("a")})
	})
	t.Run("fail: req|dive|req", func(t *testing.T) {
		errMustInclude(t, struct {
			A *string `valid:"req|dive|req"`
		}{A: ptr("")}, "required", "A")
	})
}

func TestValidateInt(t *testing.T) {
	type A struct {
		A int `valid:"req"`
	}
	t.Run("fail: required", func(t *testing.T) {
		errMustInclude(t, A{}, "required", "A")
	})
	t.Run("ok: required", func(t *testing.T) {
		errMustBeNil(t, A{A: 1})
	})
	type B struct {
		B int `valid:"min:3"`
	}
	t.Run("fail: min", func(t *testing.T) {
		errMustInclude(t, B{B: 2}, "min", "3")
	})
	t.Run("ok: min", func(t *testing.T) {
		errMustBeNil(t, B{B: 3})
	})
	type C struct {
		C int `valid:"max:3"`
	}
	t.Run("fail: max", func(t *testing.T) {
		errMustInclude(t, C{C: 4}, "max", "3")
	})
	t.Run("ok: max", func(t *testing.T) {
		errMustBeNil(t, C{C: 3})
	})
}

func TestValidateUint(t *testing.T) {
	type A struct {
		A uint `valid:"req"`
	}
	t.Run("fail: required", func(t *testing.T) {
		errMustInclude(t, A{}, "required", "A")
	})
	t.Run("ok: required", func(t *testing.T) {
		errMustBeNil(t, A{A: 1})
	})
}

func TestValidateFloat(t *testing.T) {
	type A struct {
		A float64 `valid:"req|min:3|max:5"`
	}
	t.Run("fail: required", func(t *testing.T) {
		errMustInclude(t, A{}, "required", "A")
	})
	t.Run("ok: required", func(t *testing.T) {
		errMustBeNil(t, A{A: 4})
	})
	t.Run("fail: min", func(t *testing.T) {
		errMustInclude(t, A{A: 2}, "min", "3")
	})
	t.Run("ok: min", func(t *testing.T) {
		errMustBeNil(t, A{A: 3})
	})
	t.Run("fail: max", func(t *testing.T) {
		errMustInclude(t, A{A: 6}, "max", "5")
	})
	t.Run("ok: max", func(t *testing.T) {
		errMustBeNil(t, A{A: 3})
	})
}

func TestValidateCustomRule(t *testing.T) {
	alpha := regexp.MustCompile("^[a-zA-Z]+$")
	govalid.Rule("alpha", func(v any) error {
		switch tv := v.(type) {
		case string:
			if !alpha.MatchString(tv) {
				return errors.New("must be letters only")
			}
			return nil
		default:
			return errors.New("must be used on string")
		}
	})
	type A struct {
		A string `valid:"req|alpha"`
	}
	t.Run("illegal: alpha", func(t *testing.T) {
		errMustInclude(t, struct {
			A int `valid:"req|alpha"`
		}{A: 5}, "string")
	})
	t.Run("fail: alpha", func(t *testing.T) {
		errMustInclude(t, A{A: "5"}, "letters")
	})
	t.Run("ok: alpha", func(t *testing.T) {
		errMustBeNil(t, A{A: "a"})
	})
}

func errMustInclude(t *testing.T, val any, msgs ...string) {
	t.Helper()
	err := govalid.Validate(val)
	if err == nil {
		t.Fatalf("expected non nil err; got nil")
	}
	for _, msg := range msgs {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("expected err to include %s; got %s", msg, err)
		}
	}
}

func errMustBeNil(t *testing.T, val any) {
	t.Helper()
	err := govalid.Validate(val)
	if err != nil {
		t.Fatalf("expected nil err; got %s", err)
	}
}

func ptr[T any](v T) *T {
	return &v
}
