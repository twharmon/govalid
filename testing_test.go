package govalid_test

import "testing"

func equals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("expected %v to equal %v", a, b)
	}
}

func notEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Fatalf("expected %v to not equal %v", a, b)
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
