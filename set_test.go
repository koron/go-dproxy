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

	names, err := New(v).Q("name").StringArray()
	if err != nil {
		t.Error(err)
	} else if !equalStrings(names, []string{"Bob", "Mike", "John"}) {
		t.Error("unexpected names:", names)
	}

	ages, err := New(v).Q("age").Int64Array()
	if err != nil {
		t.Error(err)
	} else if !equalInts(ages, []int64{20, 23, 22}) {
		t.Error("unexpected ages:", ages)
	}
	_ = ages
}
