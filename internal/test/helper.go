package test

import "testing"

func AssertEqualStrings(t testing.TB, expected, actual  string) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected: %q; got: %q", expected, actual)
	}
}
