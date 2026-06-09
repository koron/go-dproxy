package dproxy

// Proxy is a proxy to access a document (any).
type Proxy interface {
	// Nil returns true if the target value is nil.
	Nil() bool

	// Value returns the proxied value.  If there is no value, it returns an
	// error.
	Value() (any, error)

	// Bool returns the bool value.  If the value is not a bool, it returns an
	// error.
	Bool() (bool, error)

	// Int64 returns the int64 value.  If the value is not an int64, it returns
	// an error.
	Int64() (int64, error)

	// Float64 returns the float64 value.  If the value is not a float64, it
	// returns an error.
	Float64() (float64, error)

	// String returns the string value.  If the value is not a string, it
	// returns an error.
	String() (string, error)

	// Array returns the []any value.  If the value is not an array, it returns
	// an error.
	Array() ([]any, error)

	// Map returns the map[string]any value.  If the value is not a map, it
	// returns an error.
	Map() (map[string]any, error)

	// A returns an item from the value treated as an array.
	A(n int) Proxy

	// M returns an item from the value treated as a map.
	M(k string) Proxy

	// P returns the value pointed to by JSON Pointer query q.
	P(q string) Proxy

	// ProxySet returns a set from the value (array, slice, or map).
	ProxySet() ProxySet

	// Q returns a set of all items whose property matches k.
	Q(k string) ProxySet

	// findJPT returns a match for JSON Pointer token t.
	findJPT(t string) Proxy

	// frame provides parent frame tracking for error reporting.
	frame
}

// ProxySet provides proxy access to a set.
type ProxySet interface {
	// Empty returns true when the set is empty.
	Empty() bool

	// Len returns the count of items in the set.
	Len() int

	// BoolArray returns []bool converted from the set.
	BoolArray() ([]bool, error)

	// Int64Array returns []int64 converted from the set.
	Int64Array() ([]int64, error)

	// Float64Array returns []float64 converted from the set.
	Float64Array() ([]float64, error)

	// StringArray returns []string converted from the set.
	StringArray() ([]string, error)

	// ArrayArray returns [][]any converted from the set.
	ArrayArray() ([][]any, error)

	// MapArray returns []map[string]any converted from the set.
	MapArray() ([]map[string]any, error)

	// ProxyArray returns []Proxy that wrap each item.
	ProxyArray() ([]Proxy, error)

	// A returns a proxy for the given index in the set.
	A(n int) Proxy

	// Q returns a set of all items whose property matches k.
	Q(k string) ProxySet

	// Qc returns a set of values of property k from each item.
	Qc(k string) ProxySet

	// frame provides parent frame tracking for error reporting.
	frame
}

// New creates a new Proxy instance for v.
func New(v any) Proxy {
	return &valueProxy{value: v}
}

// NewSet creates a new ProxySet instance for v.
func NewSet(v []any) ProxySet {
	return &setProxy{values: v}
}
