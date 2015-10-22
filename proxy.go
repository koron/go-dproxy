package dproxy

// Proxy is a proxy to access a document (interface{}).
type Proxy interface {
	// Nil returns true, if target value is nil.
	Nil() bool

	// Bool returns its value.  If value isn't the type, it returns error.
	Bool() (bool, error)

	// Int64 returns its value.  If value isn't the type, it returns error.
	Int64() (int64, error)

	// Float64 returns its value.  If value isn't the type, it returns error.
	Float64() (float64, error)

	// String returns its value.  If value isn't the type, it returns error.
	String() (string, error)

	// Array returns its value.  If value isn't the type, it returns error.
	Array() ([]interface{}, error)

	// Map returns its value.  If value isn't the type, it returns error.
	Map() (map[string]interface{}, error)

	// A returns an item from value treated as the array.
	A(n int) Proxy

	// M returns an item from value treated as the map.
	M(k string) Proxy

	getParent() Proxy
	getAddress() string
}

// New creates a new proxy object.
func New(v interface{}) Proxy {
	return &valueProxy{value: v}
}
