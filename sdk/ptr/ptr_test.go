/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package ptr_test

import (
	"testing"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTo(t *testing.T) {
	p := ptr.To(1)
	require.NotNil(t, p)
	assert.Equal(t, 1, *p)
}

func TestFrom(t *testing.T) {
	t.Run("From works with int", func(t *testing.T) {
		assert.Equal(t, 0, ptr.From[int](nil))
		assert.Equal(t, 1, ptr.From(ptr.To(1)))
	})
	t.Run("From works with structx", func(t *testing.T) {
		type foo struct{ a int }
		assert.Equal(t, foo{}, ptr.From[foo](nil))
		assert.Equal(t, foo{a: 1}, ptr.From(ptr.To(foo{a: 1})))
	})
}

func TestEqual(t *testing.T) {
	t.Run("nil, nil returns true", func(t *testing.T) {
		assert.True(t, ptr.Equal[int](nil, nil))
	})
	t.Run("nil, not nil and not nil, nil returns false", func(t *testing.T) {
		assert.False(t, ptr.Equal[int](nil, ptr.To(1)))
		assert.False(t, ptr.Equal[int](ptr.To(1), nil))
	})
	t.Run("&a, &b return *a == *b", func(t *testing.T) {
		assert.False(t, ptr.Equal[int](ptr.To(1), ptr.To(2)))
		assert.True(t, ptr.Equal[int](ptr.To(1), ptr.To(1)))
	})
}
