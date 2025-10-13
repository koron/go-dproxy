package dproxy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertEqual(t *testing.T, got, want any) {
	t.Helper()
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("not equal: -want +got\n%s", d)
	}
}

func assertError(t *testing.T, want string, got error) {
	t.Helper()
	if got == nil {
		t.Fatalf("should fail with: %s", want)
	}
	if got := got.Error(); got != want {
		t.Fatalf("unexpected error:\nwant=%s\ngot=%s\n", want, got)
	}
}
