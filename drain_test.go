package dproxy

import (
	"testing"
)

func TestDrainNilReceiver(t *testing.T) {
	var d *Drain
	assertEqual(t, d.Has(), false)
	assertEqual(t, d.First(), nil)
	if got := d.All(); got != nil {
		t.Errorf("All() should return nil, got %v", got)
	}
	assertEqual(t, d.CombineErrors(), nil)
}

func TestDrainAll(t *testing.T) {
	v := parseJSON(`{"foo": 1, "bar": 2}`)

	var d Drain
	_ = d.Int64(New(v).M("foo"))
	_ = d.Int64(New(v).M("baz"))

	errs := d.All()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Error() != "not found: baz" {
		t.Errorf("unexpected error: %s", errs[0])
	}
}

func TestDrainCombineErrors(t *testing.T) {
	v := parseJSON(`{"x": 1, "y": 2, "z": 3}`)

	var d Drain
	_ = d.Int64(New(v).M("x"))
	_ = d.Int64(New(v).M("missing_a"))
	_ = d.Int64(New(v).M("missing_b"))

	err := d.CombineErrors()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertEqual(t, err.Error(), "not found: missing_a; not found: missing_b")
}

func TestDrainCombineErrorsNoError(t *testing.T) {
	v := parseJSON(`{"x": 1}`)

	var d Drain
	_ = d.Int64(New(v).M("x"))
	assertEqual(t, d.CombineErrors(), nil)
}

func TestDrainInt64(t *testing.T) {
	v := parseJSON(`{"val": 42}`)
	var d Drain
	got := d.Int64(New(v).M("val"))
	assertEqual(t, got, int64(42))
	assertEqual(t, d.Has(), false)

	got = d.Int64(New(v).M("missing"))
	assertEqual(t, got, int64(0))
	assertEqual(t, d.Has(), true)
}

func TestDrainFloat64(t *testing.T) {
	v := parseJSON(`{"val": 3.14}`)
	var d Drain
	got := d.Float64(New(v).M("val"))
	assertEqual(t, got, float64(3.14))
	assertEqual(t, d.Has(), false)
}

func TestDrainString(t *testing.T) {
	v := parseJSON(`{"val": "hello"}`)

	var d Drain
	got := d.String(New(v).M("val"))
	assertEqual(t, got, "hello")
	assertEqual(t, d.Has(), false)

	got = d.String(New(v).M("missing"))
	assertEqual(t, got, "")
	assertEqual(t, d.Has(), true)
}

func TestDrainArray(t *testing.T) {
	v := parseJSON(`{"val": ["a", "b"]}`)

	var d Drain
	got := d.Array(New(v).M("val"))
	assertEqual(t, got, []any{"a", "b"})
	assertEqual(t, d.Has(), false)
}

func TestDrainMap(t *testing.T) {
	v := parseJSON(`{"val": {"x": 1}}`)

	var d Drain
	got := d.Map(New(v).M("val"))
	assertEqual(t, got, map[string]any{"x": float64(1)})
	assertEqual(t, d.Has(), false)
}

func TestDrainBool(t *testing.T) {
	v := parseJSON(`{
		"foo": true,
		"bar": false
	}`)

	d := new(Drain)

	foo := d.Bool(New(v).M("foo"))
	if d.Has() {
		t.Error(d.First())
	} else if foo != true {
		t.Errorf("foo must be true")
	}

	bar := d.Bool(New(v).M("bar"))
	if d.Has() {
		t.Error(d.First())
	} else if bar != false {
		t.Errorf("bar must be false")
	}

	baz := d.Bool(New(v).M("baz"))
	if !d.Has() {
		t.Error("baz must not exist")
	} else if err := d.First(); err == nil || err.Error() != "not found: baz" {
		t.Errorf("unexpected error: %s", err)
	}
	_ = baz
}
