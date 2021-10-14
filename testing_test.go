package govalid_test

import (
	"strings"
	"testing"
)

func equals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("expected %v to equal %v", a, b)
	}
}

func contains(t *testing.T, haystack string, needles ...string) {
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			t.Fatalf("expected %s to contain %s", haystack, needle)
		}
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
