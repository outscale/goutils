/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package metadata_test

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/outscale/goutils/sdk/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Fetch(t *testing.T) {
	svc := metadata.NewService(http.DefaultClient)

	meta, err := svc.Fetch(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, meta.Hostname)
	assert.NotEmpty(t, meta.InstanceID)
	assert.True(t, strings.HasPrefix(meta.InstanceID, "i-"))
	assert.NotEmpty(t, meta.Placement.Subregion)
	assert.NotEmpty(t, meta.OMIID)
	assert.Truef(t, strings.HasPrefix(meta.OMIID, "ami-"), "%q should have an ami- prefix", meta.OMIID)
	assert.NotEmpty(t, meta.InstanceType)
	assert.Truef(t, strings.HasPrefix(meta.InstanceType, "tinav"), "%q should have an tinav prefix", meta.InstanceType)
	assert.NotEmpty(t, meta.MAC)
	assert.NotEmpty(t, meta.Placement.Cluster)
	assert.NotEmpty(t, meta.Placement.Server)
	assert.NotEmpty(t, meta.DeviceMapping)
	assert.NotEmpty(t, meta.Tags)
}

func TestGetHostname(t *testing.T) {
	host, err := metadata.GetHostname(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, host)
	oshost, err := os.Hostname()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(host, oshost))
}

func TestGetRegion(t *testing.T) {
	region, err := metadata.GetRegion(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, region)
	assert.Contains(t, []string{"eu-west-2", "us-east-2"}, region)
}

func TestGetSubregion(t *testing.T) {
	sr, err := metadata.GetSubregion(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, sr)
	assert.Contains(t, []string{"eu-west-2", "us-east-2"}, sr[:len(sr)-1])
	assert.Contains(t, []string{"a", "b", "c", "d"}, sr[len(sr)-1:])
}

func TestGetInstanceType(t *testing.T) {
	typ, err := metadata.GetInstanceType(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, typ)
	assert.Truef(t, strings.HasPrefix(typ, "tinav"), "%q should have an tinav prefix", typ)
}

func TestGetOMIID(t *testing.T) {
	omi, err := metadata.GetOMIID(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, omi)
	assert.Truef(t, strings.HasPrefix(omi, "ami-"), "%q should have an ami- prefix", omi)
}

func TestGetPlacementCluster(t *testing.T) {
	cluster, err := metadata.GetPlacementCluster(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, cluster)
}

func TestGetPlacementServer(t *testing.T) {
	server, err := metadata.GetPlacementServer(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, server)
}

func TestGetDeviceMapping(t *testing.T) {
	mapping, err := metadata.GetDeviceMappings(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, mapping)
	assert.NotEmpty(t, mapping["root"])
}
