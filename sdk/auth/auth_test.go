package auth_test

import (
	"testing"

	"github.com/outscale/goutils/sdk/auth"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/stretchr/testify/require"
)

func TestCheckCredentials(t *testing.T) {
	newClient := func() osc.ClientInterface {
		profile, err := profile.New(profile.FromEnv())
		if err != nil {
			return nil
		}
		client, err := osc.NewClient(profile)
		if err != nil {
			return nil
		}
		return client
	}

	t.Run("Empty credentials are rejected with an error", func(t *testing.T) {
		t.Setenv("OSC_ACCESS_KEY", "")
		t.Setenv("OSC_SECRET_KEY", "")
		t.Setenv("OSC_CONFIG_FILE", "nofile")
		t.Setenv("OSC_REGION", "eu-west-2")
		client := newClient()
		err := auth.CheckCredentials(t.Context(), client)
		require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})
	t.Run("Invalid credentials are rejected with an error", func(t *testing.T) {
		t.Setenv("OSC_ACCESS_KEY", "foo")
		t.Setenv("OSC_SECRET_KEY", "bar")
		t.Setenv("OSC_REGION", "eu-west-2")
		client := newClient()
		err := auth.CheckCredentials(t.Context(), client)
		require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})
	t.Run("CheckCredentials returns an error in case of a network problem", func(t *testing.T) {
		t.Setenv("OSC_ACCESS_KEY", "foo")
		t.Setenv("OSC_SECRET_KEY", "bar")
		t.Setenv("OSC_REGION", "foo")
		client := newClient()
		err := auth.CheckCredentials(t.Context(), client)
		require.Error(t, err)
	})
}
