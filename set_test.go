package dproxy

import "testing"

func TestSet(t *testing.T) {
	v := parseJSON(`{
		"items" : [
			{
				"name": "Bob",
				"age": 20
			},
			{
				"name": "Mike",
				"age": 23
			},
			{
				"name": "John",
				"age": 22
			}
		]
	}`)

	assertQuery(New(v).Q("name").StringArray())(
		t, []string{"Bob", "Mike", "John"})

	assertQuery(New(v).Q("age").Int64Array())(
		t, []int64{20, 23, 22})
}

func TestSetTypeError(t *testing.T) {
	v := parseJSON(`[true, false, 0, true, false]`)

	assertQerror(New(v).ProxySet().BoolArray())(
		t, "not matched types: expected=bool actual=float64: [2]")
}
