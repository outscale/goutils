/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package tags_test

import (
	"testing"

	"github.com/outscale/goutils/sdk/tags"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/stretchr/testify/assert"
)

func TestHas(t *testing.T) {
	tgs := []osc.ResourceTag{{Key: "foo", Value: "bar"}}
	tts := []struct {
		kv    []string
		found bool
	}{
		{
			found: true,
		},
		{
			kv:    []string{"foo"},
			found: true,
		},
		{
			kv:    []string{"foo", "bar"},
			found: true,
		},
		{
			kv:    []string{"bar"},
			found: false,
		},
		{
			kv:    []string{"foo", "baz"},
			found: false,
		},
		{
			kv:    []string{"bar", "baz"},
			found: false,
		},
	}
	for _, tt := range tts {
		found := tags.Has(tgs, tt.kv...)
		assert.Equal(t, tt.found, found)
	}
}

func TestHasPRefix(t *testing.T) {
	tgs := []osc.ResourceTag{{Key: "foo/bar", Value: "baz"}}
	tts := []struct {
		prefix        string
		found         bool
		suffix, value string
	}{
		{
			prefix: "foo/",
			found:  true,
			suffix: "bar",
			value:  "baz",
		},
		{
			prefix: "bar/",
			found:  false,
		},
	}
	for _, tt := range tts {
		suffix, value, found := tags.HasPrefix(tgs, tt.prefix)
		assert.Equal(t, tt.suffix, suffix)
		assert.Equal(t, tt.value, value)
		assert.Equal(t, tt.found, found)
	}
}
