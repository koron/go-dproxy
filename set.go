package dproxy

import "strconv"

type setProxy struct {
	values []interface{}
	parent frame
	label  string
}

// setProxy implements ProxySet
var _ ProxySet = (*setProxy)(nil)

func (p *setProxy) Empty() bool {
	return len(p.values) == 0
}

func (p *setProxy) Len() int {
	return len(p.values)
}

func (p *setProxy) A(n int) Proxy {
	a := "[" + strconv.Itoa(n) + "]"
	if n < 0 || n >= len(p.values) {
		return notfoundError(p, a)
	}
	return &valueProxy{
		value:  p.values[n],
		parent: p,
		label:  a,
	}
}

func (p *setProxy) Q(k string) ProxySet {
	w := findAll(p.values, k)
	return &setProxy{
		values: w,
		parent: p,
		label:  ".." + k,
	}
}

func (p *setProxy) Qc(k string) ProxySet {
	r := make([]interface{}, 0, len(p.values))
	for _, v := range p.values {
		switch v := v.(type) {
		case map[string]interface{}:
			if w, ok := v[k]; ok {
				r = append(r, w)
			}
		}
	}
	return &setProxy{
		values: r,
		parent: p,
		label:  ".." + k,
	}
}

func (p *setProxy) parentFrame() frame {
	return p.parent
}

func (p *setProxy) frameLabel() string {
	return p.label
}

func findAll(v interface{}, k string) []interface{} {
	return findAllImpl(v, k, make([]interface{}, 0, 10))
}

func findAllImpl(v interface{}, k string, r []interface{}) []interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		for n, w := range v {
			if n == k {
				r = append(r, w)
			}
			r = findAllImpl(w, k, r)
		}
	case []interface{}:
		for _, w := range v {
			r = findAllImpl(w, k, r)
		}
	}
	return r
}
