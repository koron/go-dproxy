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

func TestSet_New(t *testing.T) {
	emptySet := NewSet(nil)

	assertEqual(t, emptySet.Empty(), true)
	assertEqual(t, emptySet.Len(), 0)

	assertQuery(emptySet.BoolArray())(t, []bool{})
}

func TestSet_TypeError(t *testing.T) {
	v := parseJSON(`[true, false, 0, true, false]`)

	assertQerror(New(v).ProxySet().BoolArray())(
		t, "not matched types: expected=bool actual=float64: [2]")
}

func makeArray[T any](values ...T) []any {
	var array []any
	for _, v := range values {
		array = append(array, v)
	}
	return array
}

func TestSet_Int64Array(t *testing.T) {
	want := []int64{1, 2, 3, 4, 5}

	assertQuery(NewSet(makeArray[int](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[int32](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[int64](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[float32](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[float64](1, 2, 3, 4, 5)).Int64Array())(t, want)

	assertQerror(NewSet(makeArray[string]("foo", "bar", "baz")).Int64Array())(t, "not matched types: expected=int64 actual=string: [0]")
}

func TestSet_Float64Array(t *testing.T) {
	want := []float64{0.0, 2.0, 4.0, 6.0, 8.0}

	assertQuery(NewSet(makeArray[int](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[int32](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[int64](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[float32](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[float64](0, 2, 4, 6, 8)).Float64Array())(t, want)

	assertQerror(NewSet(makeArray[string]("foo", "bar", "baz")).Float64Array())(t, "not matched types: expected=float64 actual=string: [0]")
}

func TestSet_StringArray(t *testing.T) {
	assertQuery(NewSet(makeArray[string]("foo", "bar", "baz")).StringArray())(t, []string{"foo", "bar", "baz"})

	assertQerror(NewSet(makeArray[int64](1, 2, 3)).StringArray())(t, "not matched types: expected=string actual=int64: [0]")
}
