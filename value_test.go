package dproxy

import "testing"

func TestTypeError(t *testing.T) {
	t.Run("map at root", func(t *testing.T) {
		v := &valueProxy{}
		_, err := v.M("foo").Int64()
		assertError(t, "not required types: required=map actual=nil: (root)", err)
	})
	t.Run("map at child", func(t *testing.T) {
		v := &valueProxy{
			parent: &valueProxy{},
			label:  "foo",
		}
		_, err := v.M("bar").Int64()
		assertError(t, "not required types: required=map actual=nil: foo", err)
	})

	t.Run("array at root", func(t *testing.T) {
		v := &valueProxy{}
		_, err := v.A(0).Int64()
		assertError(t, "not required types: required=array actual=nil: (root)", err)
	})
	t.Run("array at child", func(t *testing.T) {
		v := &valueProxy{
			parent: &valueProxy{},
			label:  "foo",
		}
		_, err := v.A(0).Int64()
		assertError(t, "not required types: required=array actual=nil: foo", err)
	})
}
