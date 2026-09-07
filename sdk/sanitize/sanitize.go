package sanitize

import (
	"reflect"
	"regexp"
	"strings"
	"time"
)

type action struct {
	onField    func(reflect.StructField) bool
	valueMatch *regexp.Regexp
}

type redactFunc func(string) string

var Redacted = "[REDACTED]"

type Sanitizer struct {
	actions []action
	redact  redactFunc
}

type Option func(*Sanitizer)

// RedactAll is the default option for redaction, it replaces matched strings with [REDACTED].
func RedactAll(s *Sanitizer) {
	s.redact = func(_ string) string {
		return Redacted
	}
}

// KeepFirst2Last2 replaces matched strings with the first the first two chars + "..." + the last two chars.
func KeepFirst2Last2(s *Sanitizer) {
	s.redact = func(s string) string {
		if len(s) <= 4 {
			return Redacted
		}
		return s[:2] + "..." + s[len(s)-2:]
	}
}

// SecretKey is a regexp matching SK.
var SecretKey = regexp.MustCompile(`\b[A-Z0-9]{40}\b`)

func MatchRegexps(regexps ...*regexp.Regexp) Option {
	return func(s *Sanitizer) {
		for _, re := range regexps {
			s.actions = append(s.actions, action{
				valueMatch: re,
			})
		}
	}
}

func MatchPIIFields(s *Sanitizer) {
	s.actions = append(s.actions, action{
		onField: func(sf reflect.StructField) bool {
			return strings.Contains(sf.Tag.Get("log"), "pii")
		},
	})
}

func MatchSensitiveFields(s *Sanitizer) {
	s.actions = append(s.actions, action{
		onField: func(sf reflect.StructField) bool {
			return strings.Contains(sf.Tag.Get("log"), "sensitive")
		},
	})
}

// New will create a sanitizer that will:
// - sanitize strings (remove SK from the text using a regexp)
// - sanitice structs (redact all fields marked as pii or sensitive)
func New(opts ...Option) Sanitizer {
	if len(opts) == 0 {
		opts = []Option{MatchRegexps(SecretKey), MatchPIIFields, MatchSensitiveFields}
	}
	s := Sanitizer{}
	RedactAll(&s)
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

// Sanitize will sanitize strings, structs or slices (of strings/structx).
// A copy of the source value is returned with redacted text/fields.
func (s Sanitizer) Sanitize[T any](v T) T {
	vv := reflect.ValueOf(v)
	for _, a := range s.actions {
		vv = s.sanitizeValue(a, vv, false)
	}

	return vv.Interface().(T)
}

func (s Sanitizer) sanitizeValue(a action, v reflect.Value, withinMatch bool) reflect.Value {
	switch v.Kind() {
	case reflect.String:
		return s.sanitizeString(a, v, withinMatch)
	case reflect.Struct:
		return s.sanitizeStruct(a, v, withinMatch)
	case reflect.Slice:
		return s.sanitizeSlice(a, v, withinMatch)
	case reflect.Pointer:
		return s.sanitizePointer(a, v, withinMatch)
		// tODO: interfaces ?
	}
	return v
}

func (s Sanitizer) sanitizePointer(a action, v reflect.Value, withinMatch bool) reflect.Value {
	if v.IsNil() {
		return v
	}
	v = s.sanitizeValue(a, v.Elem(), withinMatch)
	vv := reflect.New(v.Type())
	vv.Elem().Set(v)
	return vv
}

func (s Sanitizer) sanitizeSlice(a action, v reflect.Value, withinMatch bool) reflect.Value {
	nv := reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), v.Len(), v.Len())
	for i := range v.Len() {
		nvalue := s.sanitizeValue(a, v.Index(i), withinMatch)
		nv.Index(i).Set(nvalue)
	}
	return nv
}

func (s Sanitizer) sanitizeString(a action, v reflect.Value, withinMatch bool) reflect.Value {
	if a.onField != nil && !withinMatch {
		return v
	}
	return s.doSanitizeString(a, v)
}

func (s Sanitizer) sanitizeStruct(a action, v reflect.Value, withinMatch bool) reflect.Value {
	if a.onField == nil && !withinMatch {
		return v
	}
	// time.Time constains unexported fields that cannot be read/set
	if v.Type() == reflect.TypeFor[time.Time]() {
		if withinMatch {
			return s.doSanitizeTime(v)
		}
		return v
	}
	nv := reflect.Indirect(reflect.New(v.Type()))
	for field, value := range v.Fields() {
		value = s.sanitizeValue(a, value, withinMatch || a.onField(field))
		// currently, the only non exported fields are within time.Time structs
		if field.IsExported() {
			nv.FieldByName(field.Name).Set(value)
		}
	}
	return nv
}

func (s Sanitizer) doSanitizeString(a action, v reflect.Value) reflect.Value {
	if v.IsZero() {
		return v
	}
	var nstr string
	if a.valueMatch != nil {
		nstr = a.valueMatch.ReplaceAllStringFunc(v.String(), s.redact)
	} else {
		nstr = s.redact(v.String())
	}
	return reflect.ValueOf(nstr)
}

func (s Sanitizer) doSanitizeTime(v reflect.Value) reflect.Value {
	return reflect.New(v.Type()).Elem()
}

// Default is a sanitizer with default options (strings, pii and sensitive).
var Default = New()

// Sanitize is a shortcut to Default.Sanitize
func Sanitize[T any](v T) T {
	return Default.Sanitize(v)
}
