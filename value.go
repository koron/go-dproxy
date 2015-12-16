package dproxy

import "strconv"

type valueProxy struct {
	value  interface{}
	parent frame
	label  string
}

// valueProxy implements Proxy.
var _ Proxy = (*valueProxy)(nil)

func (p *valueProxy) Nil() bool {
	return p.value == nil
}

func (p *valueProxy) Bool() (bool, error) {
	switch v := p.value.(type) {
	case bool:
		return v, nil
	default:
		return false, typeError(p, Tbool, v)
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
		return 0, typeError(p, Tint64, v)
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
		return 0, typeError(p, Tfloat64, v)
	}
}

func (p *valueProxy) String() (string, error) {
	switch v := p.value.(type) {
	case string:
		return v, nil
	default:
		return "", typeError(p, Tstring, v)
	}
}

func (p *valueProxy) Array() ([]interface{}, error) {
	switch v := p.value.(type) {
	case []interface{}:
		return v, nil
	default:
		return nil, typeError(p, Tarray, v)
	}
}

func (p *valueProxy) Map() (map[string]interface{}, error) {
	switch v := p.value.(type) {
	case map[string]interface{}:
		return v, nil
	default:
		return nil, typeError(p, Tmap, v)
	}
}

func (p *valueProxy) A(n int) Proxy {
	switch v := p.value.(type) {
	case []interface{}:
		a := "[" + strconv.Itoa(n) + "]"
		if n < 0 || n >= len(v) {
			return notfoundError(p, a)
		}
		return &valueProxy{
			value:  v[n],
			parent: p,
			label:  a,
		}
	default:
		return typeError(p, Tarray, v)
	}
}

func (p *valueProxy) M(k string) Proxy {
	switch v := p.value.(type) {
	case map[string]interface{}:
		a := "." + k
		w, ok := v[k]
		if !ok {
			return notfoundError(p, a)
		}
		return &valueProxy{
			value:  w,
			parent: p,
			label:  a,
		}
	default:
		return typeError(p, Tmap, v)
	}
}

func (p *valueProxy) ProxySet() ProxySet {
	switch v := p.value.(type) {
	case []interface{}:
		return &setProxy{
			values: v,
			parent: p,
		}
	default:
		return typeError(p, Tarray, v)
	}
}

func (p *valueProxy) Q(k string) ProxySet {
	w := findAll(p.value, k)
	return &setProxy{
		values: w,
		parent: p,
		label:  ".." + k,
	}
}

func (p *valueProxy) parentFrame() frame {
	return p.parent
}

func (p *valueProxy) frameLabel() string {
	return p.label
}
