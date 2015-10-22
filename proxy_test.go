package dproxy

import (
	"encoding/json"
	"testing"
)

func parseJson(s string) interface{} {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		panic(err)
	}
	return v
}

func TestReadme(t *testing.T) {
	v := parseJson(`{
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

	_, err = New(v).M("data").M("castom").String()
	if err == nil || err.Error() != "not found: data.castom" {
		t.Error("err is not \"not found: data.castom\":", err)
	}
}

func TestMapBool(t *testing.T) {
	v := parseJson(`{
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
