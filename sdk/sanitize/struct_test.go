package sanitize_test

import (
	"testing"

	"github.com/outscale/goutils/sdk/sanitize"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStruct(t *testing.T) {
	t.Run("Profiles are sanitized", func(t *testing.T) {
		p := profile.Profile{
			AccessKey:   "foo",
			SecretKey:   "bar",
			SecretKeyV2: "baz",
		}
		sanitized := sanitize.Struct(p)
		assert.Equal(t, sanitize.Redacted, sanitized.SecretKey)
	})
	t.Run("sensitive fields are redacted", func(t *testing.T) {
		resp := osc.CreateAccessKeyResponse{
			AccessKey: &osc.AccessKeySecretKey{
				AccessKeyId: new("foo"),
				SecretKey:   new("bar"),
			},
		}
		sanitized := sanitize.Struct(resp)
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
		sanitized := sanitize.Struct(resp)
		assert.Empty(t, *sanitized.AccessKey.SecretKey)
	})
	t.Run("pii fields are redacted", func(t *testing.T) {
		resp := osc.ReadAccountsResponse{
			Accounts: &[]osc.Account{{
				AccountId: new("foo"),
				FirstName: new("bar"),
			}},
		}
		sanitized := sanitize.Struct(resp)
		require.NotNil(t, sanitized.Accounts)
		require.Len(t, *sanitized.Accounts, 1)
		assert.NotEqual(t, sanitize.Redacted, *(*sanitized.Accounts)[0].AccountId, "non pii field must not have been redacted")
		assert.Equal(t, sanitize.Redacted, *(*sanitized.Accounts)[0].FirstName, "pii field must have been redacted")
		assert.NotEqual(t, sanitize.Redacted, *(*resp.Accounts)[0].FirstName, "source struct must not have been modified")
	})
	t.Run("Exported time fields are properly copied", func(t *testing.T) {
		{
			vm := osc.Vm{
				CreationDate: iso8601.Now(),
			}
			sanitized := sanitize.Struct(vm)
			assert.Equal(t, vm.CreationDate.String(), sanitized.CreationDate.String())
		}
		{
			vm := osc.VmGroup{
				CreationDate: new(iso8601.Now()),
			}
			sanitized := sanitize.Struct(vm)
			assert.Equal(t, vm.CreationDate.String(), sanitized.CreationDate.String())
		}
	})
}
