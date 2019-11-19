package dproxy

import (
	"fmt"
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, actual, expected interface{}, format string, a ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		msg := fmt.Sprintf(format, a...)
		t.Errorf("not equal: %s\nexpect=%+v\nactual=%+v", msg, expected, actual)
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
