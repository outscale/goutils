package sanitize_test

import (
	"strings"
	"testing"

	"github.com/outscale/goutils/sdk/sanitize"
	"github.com/stretchr/testify/assert"
)

const sk = "0A1B2C3D4D5E6F7G8H9I0A1B2C3D4D5E6F7G8H9I"

type tc struct {
	s        string
	sanitize bool
}

func testCases() []tc {
	tcs := []tc{
		{s: sk, sanitize: true},
		{s: sk + " " + sk, sanitize: true},
		{s: strings.ToLower(sk), sanitize: false},
		{s: sk[:len(sk)-1], sanitize: false},
		{s: sk + "1", sanitize: false},
	}
	for _, boundary := range []string{":", `"`, `'`, " ", "\n"} {
		tcs = append(tcs, tc{s: boundary + sk + boundary, sanitize: true})
	}
	return tcs
}

func TestString(t *testing.T) {
	for _, tc := range testCases() {
		sanitized := sanitize.String(tc.s)
		assert.Equalf(t, tc.sanitize, tc.s != sanitized, "%q: %q should have been sanitized", tc.s, sanitized)
		if tc.sanitize {
			assert.NotContainsf(t, sanitized, sk, "%q: sanitized string %q should not contain sk %q", tc.s, sanitized, sk)
		}
	}
}

func TestKeepFirst2Last2(t *testing.T) {
	s := sanitize.New(sanitize.MatchRegexps(sanitize.SecretKey), sanitize.RedactString(sanitize.KeepFirst2Last2))
	sanitized := s.Sanitize(sk)
	assert.Equal(t, "0A...9I", sanitized)
}

func TestSanitizer(t *testing.T) {
	s := sanitize.New()
	for _, tc := range testCases() {
		sanitized := s.Sanitize(tc.s)
		assert.Equalf(t, tc.sanitize, tc.s != sanitized, "%q: %q should have been sanitized", tc.s, sanitized)
		if tc.sanitize {
			assert.NotContainsf(t, sanitized, sk, "%q: sanitized string %q should not contain sk %q", tc.s, sanitized, sk)
		}
	}
}
