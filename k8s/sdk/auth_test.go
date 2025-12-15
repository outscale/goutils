package sdk_test

import (
	"testing"

	"github.com/outscale/goutils/k8s/sdk"
	"github.com/stretchr/testify/require"
)

func TestCheckCredentials(t *testing.T) {
	t.Run("Invalid credentials are rejected with an error", func(t *testing.T) {
		t.Setenv("OSC_ACCESS_KEY", "foo")
		t.Setenv("OSC_SECRET_KEY", "bar")
		_, _, err := sdk.NewSDKClient(t.Context(), "goutils/dev", sdk.Options{})
		require.ErrorIs(t, err, sdk.ErrInvalidCredentials)
	})
	t.Run("CheckCredentials returns an error in case of a network problem", func(t *testing.T) {
		t.Setenv("OSC_ACCESS_KEY", "foo")
		t.Setenv("OSC_SECRET_KEY", "bar")
		t.Setenv("OSC_REGION", "foo")
		_, _, err := sdk.NewSDKClient(t.Context(), "goutils/dev", sdk.Options{})
		require.Error(t, err)
	})
}
