package dproxy

import "strconv"

type valueProxy struct {
	value   interface{}
	parent  Proxy
	address string
}

func newValueProxy(v interface{}, parent Proxy) *valueProxy {
	return &valueProxy{
		value:  v,
		parent: parent,
	}
}

func (p *valueProxy) Nil() bool {
	return p.value == nil
}

func (p *valueProxy) Bool() (bool, error) {
	switch v := p.value.(type) {
	case bool:
		return v, nil
	default:
		return false, mismatchError(p, Tbool, v)
	}
}

func (p *valueProxy) Int64() (int64, error) {
	switch v := p.value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, mismatchError(p, Tint64, v)
	}
}

func (p *valueProxy) Float64() (float64, error) {
	switch v := p.value.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, mismatchError(p, Tfloat64, v)
	}
}

func (p *valueProxy) String() (string, error) {
	switch v := p.value.(type) {
	case string:
		return v, nil
	default:
		return "", mismatchError(p, Tstring, v)
	}
}

func (p *valueProxy) Array() ([]interface{}, error) {
	switch v := p.value.(type) {
	case []interface{}:
		return v, nil
	default:
		return nil, mismatchError(p, Tarray, v)
	}
}

func (p *valueProxy) Map() (map[string]interface{}, error) {
	switch v := p.value.(type) {
	case map[string]interface{}:
		return v, nil
	default:
		return nil, mismatchError(p, Tmap, v)
	}
}

func (p *valueProxy) A(n int) Proxy {
	switch v := p.value.(type) {
	case []interface{}:
		a := "[" + strconv.Itoa(n) + "]"
		if n < 0 || n >= len(v) {
			return addressError(p, a)
		}
		return &valueProxy{
			value:   v[n],
			parent:  p,
			address: a,
		}
	default:
		return mismatchError(p, Tarray, v)
	}
}

func (p *valueProxy) M(k string) Proxy {
	switch v := p.value.(type) {
	case map[string]interface{}:
		a := "." + k
		w, ok := v[k]
		if !ok {
			return addressError(p, a)
		}
		return &valueProxy{
			value:   w,
			parent:  p,
			address: a,
		}
	default:
		return mismatchError(p, Tmap, v)
	}
}

func (p *valueProxy) getParent() Proxy {
	return p.parent
}

func (p *valueProxy) getAddress() string {
	return p.address
}
