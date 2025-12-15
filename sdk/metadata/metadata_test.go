/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package metadata_test

import (
	"net/http"
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
	assert.NotEmpty(t, meta.InstanceID)
	assert.True(t, strings.HasPrefix(meta.InstanceID, "i-"))
	assert.NotEmpty(t, meta.Region)
	assert.Contains(t, []string{"eu-west-2", "us-east-2"}, meta.Region)
	assert.NotEmpty(t, meta.SubRegion)
	assert.True(t, strings.HasPrefix(meta.SubRegion, meta.Region))
}

func TestGetRegion(t *testing.T) {
	region, err := metadata.GetRegion(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, region)
	assert.Contains(t, []string{"eu-west-2", "us-east-2"}, region)
}
