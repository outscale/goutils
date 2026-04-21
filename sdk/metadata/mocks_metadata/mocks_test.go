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
