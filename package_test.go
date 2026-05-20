package dproxy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertEquals(t *testing.T, want, got any) {
	t.Helper()
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("not equal -want +got\n%s", d)
	}
}

func assertError(t *testing.T, err error, exp string) {
	t.Helper()
	if err == nil {
		t.Fatalf("should fail with: %s", exp)
	}
	if act := err.Error(); act != exp {
		t.Fatalf("unexpected error:\nexpect=%s\nactual=%s\n", exp, act)
	}
}
