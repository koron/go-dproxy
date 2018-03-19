package dproxy

import "testing"

func TestReflect_Readme(t *testing.T) {
	v := parseJSON(`{
		"cities": [ "tokyo", 100, "osaka", 200, "hakata", 300 ],
		"data": {
			"custom": [ "male", 21, "female", 22 ]
		}
	}`)

	s, err := NewReflect(v).M("cities").A(0).String()
	if s != "tokyo" {
		t.Error("cities[0] must be \"tokyo\":", err)
	}

	_, err = NewReflect(v).M("cities").A(0).Float64()
	if err == nil {
		t.Error("cities[0] (float64) must be failed:", err)
	}

	n, err := NewReflect(v).M("cities").A(1).Float64()
	if n != 100 {
		t.Error("cities[1] must be 100:", err)
	}

	s2, err := NewReflect(v).M("data").M("custom").A(2).String()
	if s2 != "female" {
		t.Error("data.custom[2] must be \"female\":", err)
	}

	_, err = NewReflect(v).M("data").M("kustom").String()
	if err == nil || err.Error() != "not found: data.kustom" {
		t.Error("err is not \"not found: data.kustom\":", err)
	}
}

func TestReflect_MapBool(t *testing.T) {
	v := parseJSON(`{
		"foo": true,
		"bar": false
	}`)

	// check "foo"
	foo, err := NewReflect(v).M("foo").Bool()
	if err != nil {
		t.Error(err)
	} else if foo != true {
		t.Errorf("foo must be true")
	}

	// check "bar"
	bar, err := NewReflect(v).M("bar").Bool()
	if err != nil {
		t.Error(err)
	} else if bar != false {
		t.Errorf("bar must be false")
	}
}

func TestReflectMisc(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if NewReflect(nil).Nil() == false {
			t.Error("nil.Nil() should true but false")
		}
		if NewReflect("foo").Nil() == true {
			t.Error("string.Nil() should false but true")
		}
	})
}
