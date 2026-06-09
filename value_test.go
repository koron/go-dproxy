package dproxy

import (
	"errors"
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
	assertQuery(New(any(int32(42))).Int64())(t, int64(42))
	assertQuery(New(any(float32(42))).Int64())(t, int64(42))

	assertQuery(New(customInt64er{val: 99}).Int64())(t, int64(99))

	assertQerror(New(customInt64er{err: errors.New("oops")}).Int64())(
		t, "convert error: oops: (root)")

	assertQerror(New("x").Int64())(
		t, "not matched types: expected=int64 actual=string: (root)")
}

func TestValueFloat64(t *testing.T) {
	assertQuery(New(any(int(42))).Float64())(t, float64(42))
	assertQuery(New(any(int32(42))).Float64())(t, float64(42))
	assertQuery(New(any(int64(42))).Float64())(t, float64(42))
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

func TestValueProxySet_error(t *testing.T) {
	assertQerror(New("not array").ProxySet().BoolArray())(
		t, "not matched types: expected=array actual=string: (root)")
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
