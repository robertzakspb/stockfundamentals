package test

import "testing"

func AssertEqualStrings(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("expected: %q; got: %q", want, got)
	}
}
