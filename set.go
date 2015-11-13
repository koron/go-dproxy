package dproxy

import "strconv"

type setProxy struct {
	proxies []Proxy
	parent  frame
	label   string
}

// setProxy implements ProxySet
var _ ProxySet = (*setProxy)(nil)

func (p *setProxy) Empty() bool {
	return len(p.proxies) == 0
}

func (p *setProxy) Len() int {
	return len(p.proxies)
}

func (p *setProxy) A(n int) Proxy {
	a := "[" + strconv.Itoa(n) + "]"
	if n < 0 || n >= len(p.proxies) {
		return notfoundError(p, a)
	}
	return p.proxies[n]
}

func (p *setProxy) Q(k string) ProxySet {
	// TODO: impl me
	return nil
}

func (p *setProxy) Qc(k string) ProxySet {
	// TODO: impl me
	return nil
}

func (p *setProxy) parentFrame() frame {
	return p.parent
}

func (p *setProxy) frameLabel() string {
	return p.label
}
