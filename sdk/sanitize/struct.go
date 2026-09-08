package sanitize

import (
	"reflect"
	"slices"
	"time"
)

type fieldMatcher func(reflect.StructField) bool

type StructSanitizer struct {
	matchers []fieldMatcher
	redact   redactFunc
}

type StructOption func(*StructSanitizer)

// RedactField configures the redaction function.
func RedactField(fn redactFunc) StructOption {
	return func(s *StructSanitizer) {
		s.redact = fn
	}
}

func MatchField(sensitivity ...string) StructOption {
	return func(s *StructSanitizer) {
		s.matchers = append(s.matchers, func(sf reflect.StructField) bool {
			return slices.Contains(sensitivity, sf.Tag.Get("log"))
		})
	}
}

// NewStruct will create a sanitizer that will redact all fields marked as pii or sensitive from a struct.
func NewStruct(opts ...StructOption) StructSanitizer {
	if len(opts) == 0 {
		opts = []StructOption{MatchField(PII, Sensitive)}
	}
	s := StructSanitizer{
		redact: RedactAll,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

// Sanitize will sanitize structs or slices of structs.
// A copy of the source value is returned with redacted text/fields.
func (s StructSanitizer) Sanitize[T any](v T) T {
	vv := reflect.ValueOf(v)
	for _, m := range s.matchers {
		vv = s.sanitizeValue(m, vv, false)
	}

	return vv.Interface().(T)
}

func (s StructSanitizer) sanitizeValue(m fieldMatcher, v reflect.Value, withinMatch bool) reflect.Value {
	switch v.Kind() {
	case reflect.String:
		return s.sanitizeString(m, v, withinMatch)
	case reflect.Struct:
		return s.sanitizeStruct(m, v, withinMatch)
	case reflect.Slice:
		return s.sanitizeSlice(m, v, withinMatch)
	case reflect.Pointer:
		return s.sanitizePointer(m, v, withinMatch)
		// tODO: interfaces ?
	}
	return v
}

func (s StructSanitizer) sanitizePointer(m fieldMatcher, v reflect.Value, withinMatch bool) reflect.Value {
	if v.IsNil() {
		return v
	}
	v = s.sanitizeValue(m, v.Elem(), withinMatch)
	vv := reflect.New(v.Type())
	vv.Elem().Set(v)
	return vv
}

func (s StructSanitizer) sanitizeSlice(m fieldMatcher, v reflect.Value, withinMatch bool) reflect.Value {
	nv := reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), v.Len(), v.Len())
	for i := range v.Len() {
		nvalue := s.sanitizeValue(m, v.Index(i), withinMatch)
		nv.Index(i).Set(nvalue)
	}
	return nv
}

func (s StructSanitizer) sanitizeString(m fieldMatcher, v reflect.Value, withinMatch bool) reflect.Value {
	if !withinMatch {
		return v
	}
	return s.doSanitizeString(v)
}

func (s StructSanitizer) sanitizeStruct(m fieldMatcher, v reflect.Value, withinMatch bool) reflect.Value {
	// time.Time constains unexported fields that cannot be read/set
	if v.Type() == reflect.TypeFor[time.Time]() {
		if withinMatch {
			return s.doSanitizeTime(v)
		}
		return v
	}
	nv := reflect.Indirect(reflect.New(v.Type()))
	for field, value := range v.Fields() {
		value = s.sanitizeValue(m, value, withinMatch || m(field))
		// currently, the only non exported fields are within time.Time structs
		if field.IsExported() {
			nv.FieldByName(field.Name).Set(value)
		}
	}
	return nv
}

func (s StructSanitizer) doSanitizeString(v reflect.Value) reflect.Value {
	if v.IsZero() {
		return v
	}
	return reflect.ValueOf(s.redact(v.String()))
}

func (s StructSanitizer) doSanitizeTime(v reflect.Value) reflect.Value {
	return reflect.New(v.Type()).Elem()
}

// defaultStruct is a StructSanitizer with default options (pii and sensitive).
var defaultStruct = NewStruct()

// Sanitize is a shortcut to Default.Sanitize
func Struct[T any](v T) T {
	return defaultStruct.Sanitize(v)
}
