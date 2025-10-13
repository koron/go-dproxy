package dproxy

import (
	"testing"
)

func TestReadme(t *testing.T) {
	proxy := New(parseJSON(`{
		"cities": [ "tokyo", 100, "osaka", 200, "hakata", 300 ],
		"data": {
			"custom": [ "male", 21, "female", 22 ]
		}
	}`))

	assertQuery(proxy.M("cities").A(0).String())(
		t, "tokyo")

	assertQerror(proxy.M("cities").A(0).Float64())(
		t, "not matched types: expected=float64 actual=string: cities[0]")

	assertQuery(proxy.M("cities").A(1).Float64())(
		t, 100.0)

	assertQuery(proxy.M("data").M("custom").A(2).String())(
		t, "female")

	assertQerror(proxy.M("data").M("kustom").String())(
		t, "not found: data.kustom")
}

func TestMapBool(t *testing.T) {
	v := parseJSON(`{
		"foo": true,
		"bar": false
	}`)

	// check "foo"
	assertQuery(New(v).M("foo").Bool())(t, true)

	// check "bar"
	assertQuery(New(v).M("bar").Bool())(t, false)
}

type wrappedMap map[string]any

func TestWrappedMap(t *testing.T) {
	v := wrappedMap{
		"foo": 123,
	}
	assertQuery(New(v).M("foo").Int64())(t, 123)
}
