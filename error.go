package dproxy

import "fmt"

// ErrorType is type of errors
type ErrorType int

const (
	// Etype means expected type is not matched with actual.
	Etype ErrorType = iota + 1

	// Enotfound means key or index doesn't exist.
	Enotfound
)

// Error get detail information of the errror.
type Error interface {
	// ErrorType returns type of error.
	ErrorType() ErrorType

	// FullAddress returns query string where cause first error.
	FullAddress() string
}

type errorProxy struct {
	errorType ErrorType
	parent    Proxy
	address   string
	expected  Type
	actual    Type
}

func (p *errorProxy) Nil() bool {
	return false
}

func (p *errorProxy) Bool() (bool, error) {
	return false, p
}

func (p *errorProxy) Int64() (int64, error) {
	return 0, p
}

func (p *errorProxy) Float64() (float64, error) {
	return 0, p
}

func (p *errorProxy) String() (string, error) {
	return "", p
}

func (p *errorProxy) Array() ([]interface{}, error) {
	return nil, p
}

func (p *errorProxy) Map() (map[string]interface{}, error) {
	return nil, p
}

func (p *errorProxy) A(n int) Proxy {
	return p
}

func (p *errorProxy) M(k string) Proxy {
	return p
}

func (p *errorProxy) getParent() Proxy {
	return p.parent
}

func (p *errorProxy) getAddress() string {
	return p.address
}

func (p *errorProxy) Error() string {
	switch p.errorType {
	case Etype:
		return fmt.Sprintf("not matched types: expected=%s actual=%s: %s",
			p.expected.String(), p.actual.String(), p.FullAddress())
	case Enotfound:
		return fmt.Sprintf("not found: %s", p.FullAddress())
	default:
		return fmt.Sprintf("unexpected: %s", p.FullAddress())
	}
}

func (p *errorProxy) ErrorType() ErrorType {
	return p.errorType
}

func (p *errorProxy) FullAddress() string {
	x := 0
	for q := Proxy(p); q != nil; q = q.getParent() {
		x += len(q.getAddress())
	}
	b := make([]byte, x)
	for q := Proxy(p); q != nil; q = q.getParent() {
		x -= len(q.getAddress())
		copy(b[x:], q.getAddress())
	}
	if b[0] == '.' {
		return string(b[1:])
	}
	return string(b)
}

func mismatchError(p Proxy, expected Type, actual interface{}) *errorProxy {
	return &errorProxy{
		errorType: Etype,
		parent:    p,
		expected:  expected,
		actual:    detectType(actual),
	}
}

func addressError(p Proxy, address string) *errorProxy {
	return &errorProxy{
		errorType: Enotfound,
		parent:    p,
		address:   address,
	}
}
