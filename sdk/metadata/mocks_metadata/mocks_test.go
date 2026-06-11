package mocks_metadata_test

import (
	"testing"

	"github.com/outscale/goutils/sdk/metadata"
	"github.com/outscale/goutils/sdk/metadata/mocks_metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockSubregion(t *testing.T) {
	mocks_metadata.Setup()
	defer mocks_metadata.Teardown()

	mocks_metadata.MockSubRegion("foo")

	az, err := metadata.GetSubregion(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "foo", az)
}

func TestMockInstanceID(t *testing.T) {
	mocks_metadata.Setup()
	defer mocks_metadata.Teardown()

	mocks_metadata.MockInstanceID("foo")

	az, err := metadata.GetInstanceID(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "foo", az)
}

func TestMockDeviceMapping(t *testing.T) {
	{
		mocks_metadata.Setup()
		defer mocks_metadata.Teardown()

		mocks_metadata.MockDevideMappings(nil)

		mappings, err := metadata.GetDeviceMappings(t.Context())
		require.NoError(t, err)
		assert.Equal(t, map[string]string{
			"ami":  "/dev/sda1",
			"root": "/dev/sda1",
		}, mappings)
	}
	{
		mocks_metadata.Setup()
		defer mocks_metadata.Teardown()

		mocks_metadata.MockDevideMappings(map[string]string{
			"ebs0": "/dev/xvda",
			"ebs1": "/dev/xvdb",
		})

		mappings, err := metadata.GetDeviceMappings(t.Context())
		require.NoError(t, err)
		assert.Equal(t, map[string]string{
			"ami":  "/dev/sda1",
			"root": "/dev/sda1",
			"ebs0": "/dev/xvda",
			"ebs1": "/dev/xvdb",
		}, mappings)
	}
}
