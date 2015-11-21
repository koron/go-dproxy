package dproxy

import (
	"encoding/json"
	"testing"
)

func parseJSON(s string) interface{} {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		panic(err)
	}
	return v
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		if s != b[i] {
			return false
		}
	}
	return true
}

func equalInts(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		if s != b[i] {
			return false
		}
	}
	return true
}

func TestReadme(t *testing.T) {
	v := parseJSON(`{
		"cities": [ "tokyo", 100, "osaka", 200, "hakata", 300 ],
		"data": {
			"custom": [ "male", 21, "female", 22 ]
		}
	}`)

	s, err := New(v).M("cities").A(0).String()
	if s != "tokyo" {
		t.Error("cities[0] must be \"tokyo\":", err)
	}

	_, err = New(v).M("cities").A(0).Float64()
	if err == nil {
		t.Error("cities[0] (float64) must be failed:", err)
	}

	n, err := New(v).M("cities").A(1).Float64()
	if n != 100 {
		t.Error("cities[1] must be 100:", err)
	}

	s2, err := New(v).M("data").M("custom").A(2).String()
	if s2 != "female" {
		t.Error("data.custom[2] must be \"female\":", err)
	}

	_, err = New(v).M("data").M("kustom").String()
	if err == nil || err.Error() != "not found: data.kustom" {
		t.Error("err is not \"not found: data.kustom\":", err)
	}
}

func TestMapBool(t *testing.T) {
	v := parseJSON(`{
		"foo": true,
		"bar": false
	}`)

	// check "foo"
	foo, err := New(v).M("foo").Bool()
	if err != nil {
		t.Error(err)
	} else if foo != true {
		t.Errorf("foo must be true")
	}

	// check "bar"
	bar, err := New(v).M("bar").Bool()
	if err != nil {
		t.Error(err)
	} else if bar != false {
		t.Errorf("bar must be false")
	}
}
