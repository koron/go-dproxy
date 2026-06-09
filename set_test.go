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
	assertQuery(NewSet(makeArray[int8](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[int16](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[int32](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[int64](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[uint](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[uint8](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[uint16](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[uint32](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[uint64](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[float32](1, 2, 3, 4, 5)).Int64Array())(t, want)
	assertQuery(NewSet(makeArray[float64](1, 2, 3, 4, 5)).Int64Array())(t, want)

	assertQerror(NewSet(makeArray[string]("foo", "bar", "baz")).Int64Array())(t, "not matched types: expected=int64 actual=string: [0]")
}

func TestSet_Float64Array(t *testing.T) {
	want := []float64{0.0, 2.0, 4.0, 6.0, 8.0}

	assertQuery(NewSet(makeArray[int](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[int8](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[int16](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[int32](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[int64](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[uint](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[uint8](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[uint16](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[uint32](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[uint64](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[float32](0, 2, 4, 6, 8)).Float64Array())(t, want)
	assertQuery(NewSet(makeArray[float64](0, 2, 4, 6, 8)).Float64Array())(t, want)

	assertQerror(NewSet(makeArray[string]("foo", "bar", "baz")).Float64Array())(t, "not matched types: expected=float64 actual=string: [0]")
}

func TestSet_StringArray(t *testing.T) {
	assertQuery(NewSet(makeArray[string]("foo", "bar", "baz")).StringArray())(t, []string{"foo", "bar", "baz"})

	assertQerror(NewSet(makeArray[int64](1, 2, 3)).StringArray())(t, "not matched types: expected=string actual=int64: [0]")
}

func TestSet_ArrayArray(t *testing.T) {
	v := []any{[]any{"a"}, []any{"b"}}
	assertQuery(NewSet(v).ArrayArray())(t, [][]any{{"a"}, {"b"}})

	assertQerror(NewSet(makeArray[string]("x", "y")).ArrayArray())(
		t, "not matched types: expected=array actual=string: [0]")
}

func TestSet_MapArray(t *testing.T) {
	v := []any{map[string]any{"x": int(1)}, map[string]any{"y": int(2)}}
	assertQuery(NewSet(v).MapArray())(t, []map[string]any{{"x": int(1)}, {"y": int(2)}})

	assertQerror(NewSet(makeArray[string]("x", "y")).MapArray())(
		t, "not matched types: expected=map actual=string: [0]")
}

func TestSet_ProxyArray(t *testing.T) {
	v := []any{int64(1), int64(2)}
	proxies, err := NewSet(v).ProxyArray()
	if err != nil {
		t.Fatal(err)
	}
	if len(proxies) != 2 {
		t.Fatalf("expected 2 proxies, got %d", len(proxies))
	}
	assertQuery(proxies[0].Int64())(t, int64(1))
	assertQuery(proxies[1].Int64())(t, int64(2))
}

func TestSet_A(t *testing.T) {
	v := parseJSON(`[10, 20, 30]`)

	assertQuery(New(v).ProxySet().A(0).Int64())(t, int64(10))
	assertQuery(New(v).ProxySet().A(2).Int64())(t, int64(30))

	assertQerror(New(v).ProxySet().A(-1).Int64())(
		t, "not found: [-1]")

	assertQerror(New(v).ProxySet().A(3).Int64())(
		t, "not found: [3]")
}

func TestSet_Q(t *testing.T) {
	v := parseJSON(`{
		"items": [
			{ "name": "Alice", "role": "admin" },
			{ "name": "Bob",   "role": "user" },
			{ "name": "Carol", "role": "admin" }
		]
	}`)

	assertQuery(New(v).Q("role").StringArray())(
		t, []string{"admin", "user", "admin"})
}

func TestSet_Qc(t *testing.T) {
	v := parseJSON(`[
		{ "name": "Alice", "age": 30 },
		{ "name": "Bob" },
		{ "name": "Carol", "age": 25 }
	]`)

	assertQuery(New(v).ProxySet().Qc("age").Int64Array())(
		t, []int64{30, 25})
}

func TestSet_Qc_empty(t *testing.T) {
	v := parseJSON(`[{"x": 1}, {"x": 2}]`)

	assertQuery(New(v).ProxySet().Qc("missing").Int64Array())(
		t, []int64{})
}

func TestSet_Q_onSet(t *testing.T) {
	// setProxy.Q() does deep recursive search in each element of the set
	v := parseJSON(`[
		{"nested": {"city": "NYC"}},
		{"nested": {"city": "LA"}}
	]`)

	assertQuery(New(v).ProxySet().Q("city").StringArray())(
		t, []string{"NYC", "LA"})
}
