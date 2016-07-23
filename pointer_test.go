package dproxy

import "testing"

func TestUnescapeJPT(t *testing.T) {
	f := func(d, expect string) {
		s := unescapeJPT(d)
		if s != expect {
			t.Errorf("unescapeJPT(%q) should be %q but actually %q", d, expect, s)
		}
	}
	f("foo", "foo")
	f("bar", "bar")
	f("~0", "~")
	f("foo~0bar", "foo~bar")
	f("~1", "/")
	f("foo~1bar", "foo/bar")
	f("~01", "~1")
	f("foo~01bar", "foo~1bar")
	f("~10", "/0")
}

func TestPointerInvalidQuery(t *testing.T) {
	p := Pointer(nil, "invalid")
	err, ok := p.(*errorProxy)
	if !ok {
		t.Fatalf("it should be *errorProxy but: %+v", p)
	}
	if err.errorType != EinvalidQuery {
		t.Fatalf("errorType should be EinvalidQuery but: %s", err.errorType)
	}
}

func TestPointer(t *testing.T) {
	f := func(q string, d, expected interface{}) {
		p := Pointer(d, q)
		v, err := p.Value()
		if err != nil {
			t.Errorf("Pointer:%q for %+v failed: %s", q, d, err)
			return
		}
		assertEquals(t, v, expected, "Pointer:%q for %+v", q, d)
	}

	v := parseJSON(`{
		"cities": [ "tokyo", 100, "osaka", 200, "hakata", 300 ],
		"data": {
			"custom": [ "male", 21, "female", 22 ]
		}
	}`)
	f("", v, v)
	f("/cities", v, parseJSON(`["tokyo",100,"osaka",200,"hakata",300]`))
	f("/cities/0", v, "tokyo")
	f("/cities/1", v, float64(100))
	f("/cities/2", v, "osaka")
	f("/cities/3", v, float64(200))
	f("/cities/4", v, "hakata")
	f("/cities/5", v, float64(300))
	f("/data/custom", v, parseJSON(`["male",21,"female",22]`))

	// Example from RFC6901 https://tools.ietf.org/html/rfc6901
	w := parseJSON(`{
		"foo": ["bar", "baz"],
		"": 0,
		"a/b": 1,
		"c%d": 2,
		"e^f": 3,
		"g|h": 4,
		"i\\j": 5,
		"k\"l": 6,
		" ": 7,
		"m~n": 8
	}`)
	f("", w, w)
	f("/foo", w, parseJSON(`["bar","baz"]`))
	f("/foo/0", w, "bar")
	f("/", w, float64(0))
	f("/a~1b", w, float64(1))
	f("/c%d", w, float64(2))
	f("/e^f", w, float64(3))
	f("/g|h", w, float64(4))
	f("/i\\j", w, float64(5))
	f("/k\"l", w, float64(6))
	f("/ ", w, float64(7))
	f("/m~0n", w, float64(8))
}
