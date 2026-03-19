package test

import "testing"

// TODO: Deprecate any replace with the generic method below AssertEqual[T comparable]()
func AssertEqualStrings(t testing.TB, expected, actual string) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected: %q; got: %q", expected, actual)
	}
}

func AssertEqual[T comparable](t testing.TB, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected: %v; got: %v", expected, actual)
	}
}

func AssertNotEqual[T comparable](t testing.TB, expected, actual T) {
	t.Helper()
	if expected == actual {
		t.Errorf("Values should not be equal. Expected: %v; got: %v", expected, actual)
	}
}
