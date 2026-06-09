package dproxy

import (
	"encoding/json"
	"fmt"
	"log"
)

func Example() {
	data := `{"name": "Alice", "scores": [85, 92, 78]}`
	var v any
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		log.Fatal(err)
	}
	proxy := New(v)

	name, _ := proxy.M("name").String()
	fmt.Println(name)

	score, _ := proxy.M("scores").A(1).Int64()
	fmt.Println(score)

	// Output:
	// Alice
	// 92
}

func ExampleNew() {
	v := map[string]any{
		"items": []any{"a", "b", "c"},
	}
	p := New(v)

	item, _ := p.M("items").A(2).String()
	fmt.Println(item)

	// Output:
	// c
}

func ExampleProxy_M() {
	v := parseJSON(`{"user": {"name": "Bob", "age": 30}}`)

	name, _ := New(v).M("user").M("name").String()
	fmt.Println(name)

	age, _ := New(v).M("user").M("age").Int64()
	fmt.Println(age)

	// Output:
	// Bob
	// 30
}

func ExampleProxy_A() {
	v := parseJSON(`["apple", "banana", "cherry"]`)

	p := New(v)
	for _, i := range []int{0, 1, 2} {
		s, _ := p.A(i).String()
		fmt.Println(s)
	}

	// Output:
	// apple
	// banana
	// cherry
}

func ExampleProxy_P() {
	v := parseJSON(`{
		"store": {
			"book": {"title": "Go Programming"},
			"price": 3000
		}
	}`)

	title, _ := New(v).P("/store/book/title").String()
	fmt.Println(title)

	price, _ := New(v).P("/store/price").Int64()
	fmt.Println(price)

	// Output:
	// Go Programming
	// 3000
}

func ExampleProxy_Q() {
	v := parseJSON(`{
		"users": [
			{"name": "Alice", "role": "admin"},
			{"name": "Bob",   "role": "user"},
			{"name": "Carol", "role": "admin"}
		]
	}`)

	roles, _ := New(v).Q("role").StringArray()
	for _, r := range roles {
		fmt.Println(r)
	}

	// Unordered output:
	// admin
	// user
	// admin
}

func ExampleProxy_ProxySet() {
	// From a typed slice
	v1 := []int{10, 20, 30}
	ps1 := New(v1).ProxySet()
	fmt.Println(ps1.Len())

	// From a typed map
	v2 := map[string]int{"a": 1, "b": 2}
	ps2 := New(v2).ProxySet()
	fmt.Println(ps2.Len())

	// From a JSON array
	v3 := parseJSON(`[true, false, true]`)
	ps3 := New(v3).ProxySet()
	bools, _ := ps3.BoolArray()
	fmt.Println(bools)

	// Output:
	// 3
	// 2
	// [true false true]
}

func ExampleProxySet_Qc() {
	v := parseJSON(`[
		{"name": "Alice", "age": 30},
		{"name": "Bob"},
		{"name": "Carol", "age": 25}
	]`)

	ages, _ := New(v).ProxySet().Qc("age").Int64Array()
	for _, a := range ages {
		fmt.Println(a)
	}

	// Unordered output:
	// 30
	// 25
}

func ExampleProxySet_Q() {
	v := parseJSON(`[
		{"nested": {"city": "NYC"}},
		{"nested": {"city": "LA"}}
	]`)

	cities, _ := New(v).ProxySet().Q("city").StringArray()
	for _, c := range cities {
		fmt.Println(c)
	}

	// Unordered output:
	// NYC
	// LA
}

func ExampleProxySet_genericMap() {
	// Q works with typed maps, not just map[string]any
	v := []map[string]int{{"x": 1}, {"x": 2, "y": 3}}

	xs, _ := New(v).ProxySet().Q("x").Int64Array()
	for _, x := range xs {
		fmt.Println(x)
	}

	// Unordered output:
	// 1
	// 2
}

func ExampleProxySet_StringArray() {
	v := parseJSON(`["foo", "bar", "baz"]`)

	result, _ := New(v).ProxySet().StringArray()
	fmt.Println(result)

	// Output:
	// [foo bar baz]
}

func ExamplePointer() {
	v := parseJSON(`{"cities": ["Tokyo", "Osaka", "Hakata"]}`)

	city, _ := Pointer(v, "/cities/0").String()
	fmt.Println(city)

	// Output:
	// Tokyo
}

func ExampleDrain() {
	v := parseJSON(`{
		"a": 1,
		"b": "hello",
		"c": 3.14
	}`)

	var d Drain
	p := New(v)

	_ = d.Int64(p.M("a"))
	_ = d.String(p.M("b"))
	_ = d.Float64(p.M("c"))

	if d.Has() {
		fmt.Println("errors:", d.CombineErrors())
	} else {
		fmt.Println("all ok")
	}

	// Output:
	// all ok
}

func ExampleDrain_error() {
	v := parseJSON(`{"x": "not a number"}`)

	var d Drain
	p := New(v)

	val := d.Int64(p.M("x"))
	fmt.Println("got:", val)

	if err := d.CombineErrors(); err != nil {
		fmt.Println("error:", err)
	}

	// Output:
	// got: 0
	// error: not matched types: expected=int64 actual=string: x
}
