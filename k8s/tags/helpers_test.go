/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package tags_test

import (
	"testing"

	"github.com/outscale/goutils/k8s/tags"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/stretchr/testify/assert"
)

func TestGetClusterID(t *testing.T) {
	tgs := []osc.ResourceTag{
		{Key: "foo", Value: "bar"},
		{Key: "OscK8sClusterID/37286b54-bb4b-46c8-ac8a-37621f1e7123", Value: "owned"},
	}
	cid := tags.GetClusterID(tgs)
	assert.Equal(t, "37286b54-bb4b-46c8-ac8a-37621f1e7123", cid)
}

func TestImports(t *testing.T) {
	tgs := []osc.ResourceTag{
		{Key: "foo", Value: "bar"},
		{Key: "OscK8sClusterID/37286b54-bb4b-46c8-ac8a-37621f1e7123", Value: "owned"},
	}
	assert.True(t, tags.Has(tgs, "foo"))
	assert.Equal(t, "bar", tags.Must(tags.GetValue(tgs, "foo")))
}
