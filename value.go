package dproxy

import (
	"reflect"
	"strconv"
)

type valueProxy struct {
	value  any
	parent frame
	label  string
}

// valueProxy implements Proxy.
var _ Proxy = (*valueProxy)(nil)

func (p *valueProxy) Nil() bool {
	return p.value == nil
}

func (p *valueProxy) Value() (any, error) {
	return p.value, nil
}

func (p *valueProxy) Bool() (bool, error) {
	switch v := p.value.(type) {
	case bool:
		return v, nil
	default:
		return false, typeError(p, Tbool, v)
	}
}

type int64er interface {
	Int64() (int64, error)
}

func (p *valueProxy) Int64() (int64, error) {
	switch v := p.value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case int64er:
		w, err := v.Int64()
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
	default:
		return 0, typeError(p, Tint64, v)
	}
}

type float64er interface {
	Float64() (float64, error)
}

func (p *valueProxy) Float64() (float64, error) {
	switch v := p.value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case float64er:
		w, err := v.Float64()
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
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

func (p *valueProxy) Array() ([]any, error) {
	switch v := p.value.(type) {
	case []any:
		return v, nil
	default:
		return nil, typeError(p, Tarray, v)
	}
}

func (p *valueProxy) Map() (map[string]any, error) {
	switch v := p.value.(type) {
	case map[string]any:
		return v, nil
	default:
		return nil, typeError(p, Tmap, v)
	}
}

func (p *valueProxy) A(n int) Proxy {
	switch v := p.value.(type) {
	case []any:
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
		return requiredTypeError(p, Tarray, v)
	}
}

var mapType = reflect.TypeOf(map[string]any(nil))

func (p *valueProxy) m(v map[string]any, k string) Proxy {
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
}

func (p *valueProxy) M(k string) Proxy {
	if v, ok := p.value.(map[string]any); ok {
		return p.m(v, k)
	}

	if rv := reflect.ValueOf(p.value); rv.IsValid() && rv.Type().ConvertibleTo(mapType) {
		v, _ := rv.Convert(mapType).Interface().(map[string]any)
		return p.m(v, k)
	}

	return requiredTypeError(p, Tmap, p.value)
}

func (p *valueProxy) P(q string) Proxy {
	return pointer(p, q)
}

func (p *valueProxy) ProxySet() ProxySet {
	switch v := p.value.(type) {
	case []any:
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

func (p *valueProxy) findJPT(t string) Proxy {
	switch v := p.value.(type) {
	case map[string]any:
		return p.M(t)
	case []any:
		n, err := strconv.ParseUint(t, 10, 0)
		if err != nil {
			return &errorProxy{
				errorType: EinvalidIndex,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return p.A(int(n))
	default:
		return &errorProxy{
			errorType: EmapNorArray,
			parent:    p,
			actual:    detectType(v),
		}
	}
}

func (p *valueProxy) parentFrame() frame {
	return p.parent
}

func (p *valueProxy) frameLabel() string {
	return p.label
}
