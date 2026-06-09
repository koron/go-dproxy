package dproxy

import "testing"

func TestError_Etype_message(t *testing.T) {
	err := typeError(&simpleFrame{}, Tint64, "bad")
	assertEqual(t, err.Error(), "not matched types: expected=int64 actual=string: (root)")
}

func TestError_Enotfound_message(t *testing.T) {
	err := notfoundError(&simpleFrame{label: ".data"}, ".foo")
	assertEqual(t, err.Error(), "not found: data.foo")
}

func TestError_EmapNorArray(t *testing.T) {
	// JSON Pointer on a scalar value triggers EmapNorArray
	p := New("string value").P("/foo")
	_, err := p.String()
	if err == nil {
		t.Fatal("expected error")
	}
	assertEqual(t, err.Error(), "not map nor array: actual=string: (root)")
}

func TestError_EconvertFailure(t *testing.T) {
	p := New(customInt64er{err: errSentinel("conversion bomb")})
	_, err := p.Int64()
	if err == nil {
		t.Fatal("expected error")
	}
	assertEqual(t, err.Error(), "convert error: conversion bomb: (root)")
}

func TestError_EinvalidIndex(t *testing.T) {
	// JSON Pointer with non-numeric token on array triggers EinvalidIndex
	p := New([]any{"a", "b"}).P("/abc")
	_, err := p.String()
	if err == nil {
		t.Fatal("expected error")
	}
	assertEqual(t, err.Error(), "invalid index: strconv.ParseUint: parsing \"abc\": invalid syntax: (root)")
}

func TestError_EinvalidQuery_message(t *testing.T) {
	p := Pointer(nil, "no-slash")
	_, err := p.Value()
	if err == nil {
		t.Fatal("expected error")
	}
	assertEqual(t, err.Error(), "invalid query: not start with '/': (root)")
}

func TestError_ErequiredType_message(t *testing.T) {
	err := requiredTypeError(&simpleFrame{}, Tmap, "bad")
	assertEqual(t, err.Error(), "not required types: required=map actual=string: (root)")
}

func TestError_FullAddress(t *testing.T) {
	v := &valueProxy{
		parent: &valueProxy{label: ".data"},
		label:  ".custom",
	}
	err := typeError(v, Tstring, int64(0))
	assertEqual(t, err.Error(), "not matched types: expected=string actual=int64: data.custom")
}

func TestError_ErrorType_getter(t *testing.T) {
	err := typeError(&simpleFrame{}, Tstring, nil)
	assertEqual(t, err.ErrorType(), Etype)

	err2 := notfoundError(&simpleFrame{}, ".x")
	assertEqual(t, err2.ErrorType(), Enotfound)
}

// errSentinel is a simple error type for testing.
type errSentinel string

func (e errSentinel) Error() string { return string(e) }
