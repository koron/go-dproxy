package dproxy

import (
	"fmt"
	"reflect"
	"strconv"
)

// reflectProxy is a proxy using reflect.
type reflectProxy struct {
	rv reflect.Value
	p  frame  // parent
	l  string // label
}

var _ Proxy = (*reflectProxy)(nil)

func newReflectProxy(v interface{}, parent frame, label string) Proxy {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		if !rv.IsValid() {
			return &errorProxy{
				errorType: EinvalidValue,
				parent:    parent,
				label:     label,
				infoStr:   fmt.Sprintf("%T", v),
			}
		}
	}
	return &reflectProxy{rv: rv, p: parent, l: label}
}

func (rp *reflectProxy) parentFrame() frame {
	return rp.p
}

func (rp *reflectProxy) frameLabel() string {
	return rp.l
}

func (rp *reflectProxy) typeError(expected Type) *errorProxy {
	return typeError(rp, expected, rp.rv.Interface())
}

func (rp *reflectProxy) Nil() bool {
	return !rp.rv.IsValid()
}

func (rp *reflectProxy) Value() (interface{}, error) {
	return rp.rv.Interface(), nil
}

func (rp *reflectProxy) Bool() (bool, error) {
	switch rp.rv.Kind() {
	case reflect.Bool:
		return rp.rv.Bool(), nil
	default:
		return false, rp.typeError(Tbool)
	}
}

func (rp *reflectProxy) isInt() bool {
	switch rp.rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func (rp *reflectProxy) Int64() (int64, error) {
	if !rp.isInt() {
		return 0, rp.typeError(Tint64)
	}
	return rp.rv.Int(), nil
}

func (rp *reflectProxy) isFloat() bool {
	switch rp.rv.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func (rp *reflectProxy) Float64() (float64, error) {
	if !rp.isFloat() {
		return 0, rp.typeError(Tfloat64)
	}
	return rp.rv.Float(), nil
}

func (rp *reflectProxy) String() (string, error) {
	if rp.rv.Kind() != reflect.String {
		return "", rp.typeError(Tstring)
	}
	return rp.rv.String(), nil
}

func (rp *reflectProxy) Array() ([]interface{}, error) {
	if rp.rv.Kind() != reflect.Array {
		return nil, rp.typeError(Tarray)
	}
	// TODO: return value as []interface{}
	return nil, nil
}

func (rp *reflectProxy) Map() (map[string]interface{}, error) {
	if rp.rv.Kind() != reflect.Map {
		return nil, rp.typeError(Tmap)
	}
	// TODO: return value as map[string]interface{}
	return nil, nil
}

func (rp *reflectProxy) A(n int) Proxy {
	if rp.rv.Kind() != reflect.Array {
		return rp.typeError(Tarray)
	}
	adrs := "[" + strconv.Itoa(n) + "]"
	if n < 0 || n >= rp.rv.Len() {
		return notfoundError(rp, adrs)
	}
	v := rp.rv.Index(n)
	return newReflectProxy(v.Interface(), rp, adrs)
}

func (rp *reflectProxy) M(k string) Proxy {
	adrs := "." + k
	switch rp.rv.Kind() {
	case reflect.Map:
		v := rp.rv.MapIndex(reflect.ValueOf(k))
		if !v.IsValid() {
			return notfoundError(rp, adrs)
		}
		return newReflectProxy(v.Interface(), rp, adrs)
	case reflect.Struct:
		v := rp.rv.FieldByName(k)
		if !v.IsValid() {
			return notfoundError(rp, adrs)
		}
		return newReflectProxy(v.Interface(), rp, adrs)
	default:
		return rp.typeError(Tmap)
	}
}

func (rp *reflectProxy) P(q string) Proxy {
	return pointer(rp, q)
}

func (rp *reflectProxy) ProxySet() ProxySet {
	// TODO: return proxy set for reflect vaue.
	return nil
}

func (rp *reflectProxy) Q(k string) ProxySet {
	// TODO: return proxy set for queried value
	return nil
}

func (rp *reflectProxy) findJPT(t string) Proxy {
	switch rp.rv.Kind() {
	case reflect.Map, reflect.Struct:
		return rp.M(t)
	case reflect.Array:
		n, err := strconv.ParseUint(t, 10, 0)
		if err != nil {
			return &errorProxy{
				errorType: EinvalidIndex,
				parent:    rp,
				infoStr:   err.Error(),
			}
		}
		return rp.A(int(n))
	default:
		return &errorProxy{
			errorType: EmapNorArray,
			parent:    rp,
			actual:    detectType(rp.rv.Interface()),
		}
	}
}
