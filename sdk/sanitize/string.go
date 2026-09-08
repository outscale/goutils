package sanitize

import (
	"regexp"
)

type StringSanitizer struct {
	regexps []*regexp.Regexp
	redact  redactFunc
}

type StringOption func(*StringSanitizer)

// SecretKey is a regexp matching SK.
var SecretKey = regexp.MustCompile(`\b[A-Z0-9]{40}\b`)

// RedactString configures the redaction function.
func RedactString(fn redactFunc) StringOption {
	return func(s *StringSanitizer) {
		s.redact = fn
	}
}

func MatchRegexps(regexps ...*regexp.Regexp) StringOption {
	return func(s *StringSanitizer) {
		s.regexps = regexps
	}
}

// NewStruct will create a sanitizer that will redact matching regexps from a string.
func New(opts ...StringOption) StringSanitizer {
	if len(opts) == 0 {
		opts = []StringOption{MatchRegexps(SecretKey)}
	}
	s := StringSanitizer{
		redact: RedactAll,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

// Sanitize will sanitize strings, structs or slices (of strings/structx).
// A copy of the source value is returned with redacted text/fields.
func (s StringSanitizer) Sanitize(str string) string {
	for _, re := range s.regexps {
		str = re.ReplaceAllStringFunc(str, s.redact)
	}
	return str
}

// defaultString is a sanitizer with default options (sk).
var defaultString = New()

// Sanitize is a shortcut to Default.Sanitize
func String(s string) string {
	return defaultString.Sanitize(s)
}
