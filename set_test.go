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
		t.Fatal(err)
	} else if !equalStrings(names, []string{"Bob", "Mike", "John"}) {
		t.Error("unexpected names:", names)
	}

	ages, err := New(v).Q("age").Int64Array()
	if err != nil {
		t.Fatal(err)
	} else if !equalInts(ages, []int64{20, 23, 22}) {
		t.Error("unexpected ages:", ages)
	}
	_ = ages
}

func TestSetTypeError(t *testing.T) {
	v := parseJSON(`[true, false, 0, true, false]`)
	_, err := New(v).ProxySet().BoolArray()
	if err == nil {
		t.Fatal("should fail")
	}
	err2, ok := err.(Error)
	if !ok {
		t.Fatal("err is not Error:", err)
	}
	if et := err2.ErrorType(); et != Etype {
		t.Fatal("unexpected ErrorType:", et)
	}
	if ea := err2.FullAddress(); ea != "[2]" {
		t.Fatal("unexpected FullAddress:", ea)
	}
}
