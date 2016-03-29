package dproxy

import "bytes"

// Drain stores errors from Proxy or ProxySet.
type Drain struct {
	errors []error
}

// Has returns true if the drain stored some of errors.
func (s *Drain) Has() bool {
	return s != nil && s.errors != nil
}

// First returns a stored error.  Returns nil if there are no errors.
func (s *Drain) First() error {
	if s == nil || s.errors == nil {
		return nil
	}
	return s.errors[0]
}

// All returns all errors which stored.  Return nil if no errors stored.
func (s *Drain) All() []error {
	if s == nil || s.errors == nil {
		return nil
	}
	a := make([]error, 0, len(s.errors))
	return append(a, s.errors...)
}

// CombineErrors returns an error which combined all stored errors.  Return nil
// if not erros stored.
func (s *Drain) CombineErrors() error {
	if s == nil || s.errors == nil {
		return nil
	}
	return drainError(s.errors)
}

func (s *Drain) put(err error) {
	if err == nil {
		return
	}
	s.errors = append(s.errors, err)
}

// Bool returns bool value and stores an error.
func (s *Drain) Bool(v bool, err error) bool {
	s.put(err)
	return v
}

// Int64 returns int64 value and stores an error.
func (s *Drain) Int64(v int64, err error) int64 {
	s.put(err)
	return v
}

// Float64 returns float64 value and stores an error.
func (s *Drain) Float64(v float64, err error) float64 {
	s.put(err)
	return v
}

// String returns string value and stores an error.
func (s *Drain) String(v string, err error) string {
	s.put(err)
	return v
}

// Array returns []interface{} value and stores an error.
func (s *Drain) Array(v []interface{}, err error) []interface{} {
	s.put(err)
	return v
}

// Map returns map[string]interface{} value and stores an error.
func (s *Drain) Map(v map[string]interface{}, err error) map[string]interface{} {
	s.put(err)
	return v
}

// BoolArray returns []bool value and stores an error.
func (s *Drain) BoolArray(v []bool, err error) []bool {
	s.put(err)
	return v
}

// Int64Array returns []int64 value and stores an error.
func (s *Drain) Int64Array(v []int64, err error) []int64 {
	s.put(err)
	return v
}

// Float64Array returns []float64 value and stores an error.
func (s *Drain) Float64Array(v []float64, err error) []float64 {
	s.put(err)
	return v
}

// StringArray returns []string value and stores an error.
func (s *Drain) StringArray(v []string, err error) []string {
	s.put(err)
	return v
}

// ArrayArray returns [][]interface{} value and stores an error.
func (s *Drain) ArrayArray(v [][]interface{}, err error) [][]interface{} {
	s.put(err)
	return v
}

// MapArray returns []map[string]interface{} value and stores an error.
func (s *Drain) MapArray(v []map[string]interface{}, err error) []map[string]interface{} {
	s.put(err)
	return v
}

type drainError []error

func (derr drainError) Error() string {
	b := bytes.Buffer{}
	for i, err := range derr {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}
