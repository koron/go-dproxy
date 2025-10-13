package dproxy

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func parseJSON(s string) any {
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		panic(err)
	}
	return v
}

func assertQuery[T any](got T, err error) func(*testing.T, T) {
	return func(t *testing.T, want T) {
		t.Helper()
		if err != nil {
			t.Errorf("query failed: %s", err)
			return
		}
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("unmatched query results: -want +got\n%s", d)
		}
	}
}

func assertQerror[T any](got T, gotErr error) func(*testing.T, string) {
	return func(t *testing.T, want string) {
		t.Helper()
		if gotErr == nil {
			t.Errorf("query succeeded unexpectedly with return value: %v", got)
			return
		}
		if d := cmp.Diff(want, gotErr.Error()); d != "" {
			t.Errorf("unexpected failure: -want +got\n%s", d)
		}
	}
}
