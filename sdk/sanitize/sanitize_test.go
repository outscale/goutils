package sanitize_test

import (
	"strings"
	"testing"

	"github.com/outscale/goutils/sdk/sanitize"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		sanitized := sanitize.Sanitize(tc.s)
		assert.Equalf(t, tc.sanitize, tc.s != sanitized, "%q: %q should have been sanitized", tc.s, sanitized)
		if tc.sanitize {
			assert.NotContainsf(t, sanitized, sk, "%q: sanitized string %q should not contain sk %q", tc.s, sanitized, sk)
		}
	}
}

func TestKeepFirst2Last2(t *testing.T) {
	s := sanitize.New(sanitize.MatchRegexps(sanitize.SecretKey), sanitize.KeepFirst2Last2)
	sanitized := s.Sanitize(sk)
	assert.Equal(t, "0A...9I", sanitized)
}

func TestStringSlice(t *testing.T) {
	for _, tc := range testCases() {
		sanitized := sanitize.Sanitize([]string{tc.s, tc.s, "foo"})
		require.Len(t, sanitized, 3)
		assert.Equalf(t, tc.sanitize, tc.s != sanitized[0], "%q: %q should have been sanitized", tc.s, sanitized)
		if tc.sanitize {
			assert.NotContainsf(t, sanitized[0], sk, "%q: sanitized string %q should not contain sk %q", tc.s, sanitized, sk)
		}
		assert.Equal(t, "foo", sanitized[2])
	}
}

func TestStruct(t *testing.T) {
	t.Run("Profiles are sanitized", func(t *testing.T) {
		p := profile.Profile{
			AccessKey:   "foo",
			SecretKey:   "bar",
			SecretKeyV2: "baz",
		}
		sanitized := sanitize.Sanitize(p)
		assert.Equal(t, sanitize.Redacted, sanitized.SecretKey)
	})
	t.Run("sensitive fields are redacted", func(t *testing.T) {
		resp := osc.CreateAccessKeyResponse{
			AccessKey: &osc.AccessKeySecretKey{
				AccessKeyId: new("foo"),
				SecretKey:   new("bar"),
			},
		}
		sanitized := sanitize.Sanitize(resp)
		assert.NotEqual(t, sanitize.Redacted, *sanitized.AccessKey.AccessKeyId, "non sentitive field must not have been redacted")
		assert.Equal(t, sanitize.Redacted, *sanitized.AccessKey.SecretKey, "sentitive field must have been redacted")
		assert.NotEqual(t, sanitize.Redacted, *resp.AccessKey.SecretKey, "source struct must not have been modified")
	})
	t.Run("empty fields are not redacted", func(t *testing.T) {
		resp := osc.CreateAccessKeyResponse{
			AccessKey: &osc.AccessKeySecretKey{
				SecretKey: new(""),
			},
		}
		sanitized := sanitize.Sanitize(resp)
		assert.Empty(t, *sanitized.AccessKey.SecretKey)
	})
	t.Run("pii fields are redacted", func(t *testing.T) {
		resp := osc.ReadAccountsResponse{
			Accounts: &[]osc.Account{{
				AccountId: new("foo"),
				FirstName: new("bar"),
			}},
		}
		sanitized := sanitize.Sanitize(resp)
		require.NotNil(t, sanitized.Accounts)
		require.Len(t, *sanitized.Accounts, 1)
		assert.NotEqual(t, sanitize.Redacted, *(*sanitized.Accounts)[0].AccountId, "non pii field must not have been redacted")
		assert.Equal(t, sanitize.Redacted, *(*sanitized.Accounts)[0].FirstName, "pii field must have been redacted")
		assert.NotEqual(t, sanitize.Redacted, *(*resp.Accounts)[0].FirstName, "source struct must not have been modified")
	})
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
