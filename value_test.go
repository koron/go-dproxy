package dproxy

import (
	"errors"
	"sort"
	"testing"
)

type customInt64er struct {
	val int64
	err error
}

func (c customInt64er) Int64() (int64, error) {
	return c.val, c.err
}

type customFloat64er struct {
	val float64
	err error
}

func (c customFloat64er) Float64() (float64, error) {
	return c.val, c.err
}

func TestValueNil(t *testing.T) {
	assertEqual(t, New(nil).Nil(), true)
	assertEqual(t, New("foo").Nil(), false)
}

func TestValueArray(t *testing.T) {
	v := []any{"a", "b"}
	assertQuery(New(v).Array())(t, v)

	assertQerror(New("not an array").Array())(
		t, "not matched types: expected=array actual=string: (root)")
}

func TestValueMap(t *testing.T) {
	v := map[string]any{"foo": "bar"}
	assertQuery(New(v).Map())(t, v)

	assertQerror(New("not a map").Map())(
		t, "not matched types: expected=map actual=string: (root)")
}

func TestValueInt64(t *testing.T) {
	assertQuery(New(any(int(42))).Int64())(t, int64(42))
	assertQuery(New(any(int8(42))).Int64())(t, int64(42))
	assertQuery(New(any(int16(42))).Int64())(t, int64(42))
	assertQuery(New(any(int32(42))).Int64())(t, int64(42))
	assertQuery(New(any(uint(42))).Int64())(t, int64(42))
	assertQuery(New(any(uint8(42))).Int64())(t, int64(42))
	assertQuery(New(any(uint16(42))).Int64())(t, int64(42))
	assertQuery(New(any(uint32(42))).Int64())(t, int64(42))
	assertQuery(New(any(uint64(42))).Int64())(t, int64(42))
	assertQuery(New(any(float32(42))).Int64())(t, int64(42))

	assertQuery(New(customInt64er{val: 99}).Int64())(t, int64(99))

	assertQerror(New(customInt64er{err: errors.New("oops")}).Int64())(
		t, "convert error: oops: (root)")

	assertQerror(New("x").Int64())(
		t, "not matched types: expected=int64 actual=string: (root)")
}

func TestValueFloat64(t *testing.T) {
	assertQuery(New(any(int(42))).Float64())(t, float64(42))
	assertQuery(New(any(int8(42))).Float64())(t, float64(42))
	assertQuery(New(any(int16(42))).Float64())(t, float64(42))
	assertQuery(New(any(int32(42))).Float64())(t, float64(42))
	assertQuery(New(any(int64(42))).Float64())(t, float64(42))
	assertQuery(New(any(uint(42))).Float64())(t, float64(42))
	assertQuery(New(any(uint8(42))).Float64())(t, float64(42))
	assertQuery(New(any(uint16(42))).Float64())(t, float64(42))
	assertQuery(New(any(uint32(42))).Float64())(t, float64(42))
	assertQuery(New(any(uint64(42))).Float64())(t, float64(42))
	assertQuery(New(any(float32(1.5))).Float64())(t, float64(1.5))

	assertQuery(New(customFloat64er{val: 3.14}).Float64())(t, float64(3.14))

	assertQerror(New(customFloat64er{err: errors.New("oops")}).Float64())(
		t, "convert error: oops: (root)")

	assertQerror(New("x").Float64())(
		t, "not matched types: expected=float64 actual=string: (root)")
}

func TestValueBool_error(t *testing.T) {
	assertQerror(New("not bool").Bool())(
		t, "not matched types: expected=bool actual=string: (root)")
}

func TestValueString_error(t *testing.T) {
	assertQerror(New(42).String())(
		t, "not matched types: expected=string actual=int64: (root)")
}

func TestValueProxySet_string(t *testing.T) {
	assertQerror(New("not container").ProxySet().BoolArray())(
		t, "convert error: string is not supported for set: (root)")
}

func TestValueProxySet_map(t *testing.T) {
	v := map[string]any{"a": int64(1), "b": int64(2), "c": int64(3)}
	ps := New(v).ProxySet()
	assertEqual(t, ps.Len(), 3)

	got, err := ps.Int64Array()
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	assertEqual(t, got, []int64{1, 2, 3})
}

type wrappedMap2 map[string]any

func TestValueProxySet_wrappedMap(t *testing.T) {
	v := wrappedMap2{"x": int64(10), "y": int64(20)}
	ps := New(v).ProxySet()
	assertEqual(t, ps.Len(), 2)

	got, err := ps.Int64Array()
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	assertEqual(t, got, []int64{10, 20})
}

func TestValueProxySet_genericMap(t *testing.T) {
	v := map[string]int{"a": 1, "b": 2, "c": 3}
	ps := New(v).ProxySet()
	assertEqual(t, ps.Len(), 3)

	got, err := ps.Int64Array()
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	assertEqual(t, got, []int64{1, 2, 3})
}

func TestValueProxySet_intKeyMap(t *testing.T) {
	v := map[int]string{1: "one", 2: "two"}
	ps := New(v).ProxySet()
	assertEqual(t, ps.Len(), 2)

	got, err := ps.StringArray()
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	assertEqual(t, got, []string{"one", "two"})
}

func TestValueProxySet_intSlice(t *testing.T) {
	v := []int{10, 20, 30}
	ps := New(v).ProxySet()
	assertEqual(t, ps.Len(), 3)

	got, err := ps.Int64Array()
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, got, []int64{10, 20, 30})
}

func TestValueProxySet_stringSlice(t *testing.T) {
	v := []string{"a", "b", "c"}
	ps := New(v).ProxySet()
	assertEqual(t, ps.Len(), 3)

	got, err := ps.StringArray()
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, got, []string{"a", "b", "c"})
}

func TestValueProxySet_float64Slice(t *testing.T) {
	v := []float64{1.5, 2.5}
	ps := New(v).ProxySet()
	assertEqual(t, ps.Len(), 2)

	got, err := ps.Float64Array()
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, got, []float64{1.5, 2.5})
}

func TestValueA_outOfBounds(t *testing.T) {
	v := parseJSON(`["a", "b"]`)

	assertQerror(New(v).A(-1).String())(
		t, "not found: [-1]")

	assertQerror(New(v).A(2).String())(
		t, "not found: [2]")
}

func TestTypeError(t *testing.T) {
	t.Run("map at root", func(t *testing.T) {
		v := &valueProxy{}
		assertQerror(v.M("foo").Int64())(t, "not required types: required=map actual=nil: (root)")
	})
	t.Run("map at child", func(t *testing.T) {
		v := &valueProxy{
			parent: &valueProxy{},
			label:  "foo",
		}
		assertQerror(v.M("foo").Int64())(t, "not required types: required=map actual=nil: foo")
	})

	t.Run("array at root", func(t *testing.T) {
		v := &valueProxy{}
		assertQerror(v.A(0).Int64())(t, "not required types: required=array actual=nil: (root)")
	})
	t.Run("array at child", func(t *testing.T) {
		v := &valueProxy{
			parent: &valueProxy{},
			label:  "foo",
		}
		assertQerror(v.A(0).Int64())(t, "not required types: required=array actual=nil: foo")
	})
}
