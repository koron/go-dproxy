package dproxy

import (
	"testing"
)

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
