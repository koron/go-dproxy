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
