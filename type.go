package dproxy

// Type is the type of a value.
type Type int

const (
	// Tunknown indicates the value is not supported.
	Tunknown Type = iota

	// Tnil indicates the value is nil.
	Tnil

	// Tbool indicates the value is a bool.
	Tbool

	// Tint64 indicates the value is an int64.
	Tint64

	// Tfloat64 indicates the value is a float64.
	Tfloat64

	// Tstring indicates the value is a string.
	Tstring

	// Tarray indicates the value is an array ([]any).
	Tarray

	// Tmap indicates the value is a map (map[string]any).
	Tmap
)

// detectType returns the type of a value.
func detectType(v any) Type {
	if v == nil {
		return Tnil
	}
	switch v.(type) {
	case bool:
		return Tbool
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		return Tint64
	case float32, float64:
		return Tfloat64
	case string:
		return Tstring
	case []any:
		return Tarray
	case map[string]any:
		return Tmap
	default:
		return Tunknown
	}
}

func (t Type) String() string {
	switch t {
	case Tunknown:
		return "unknown"
	case Tnil:
		return "nil"
	case Tbool:
		return "bool"
	case Tint64:
		return "int64"
	case Tfloat64:
		return "float64"
	case Tstring:
		return "string"
	case Tarray:
		return "array"
	case Tmap:
		return "map"
	default:
		return "unknown"
	}
}
