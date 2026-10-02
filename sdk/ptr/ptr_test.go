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

func TestFrom(t *testing.T) {
	t.Run("From works with int", func(t *testing.T) {
		assert.Equal(t, 0, ptr.From[int](nil))
		assert.Equal(t, 2, ptr.From[int](nil, 2))
		assert.Equal(t, 1, ptr.From(new(1)))
	})
	t.Run("From works with structs", func(t *testing.T) {
		type foo struct{ a int }
		assert.Equal(t, foo{}, ptr.From[foo](nil))
		assert.Equal(t, foo{a: 1}, ptr.From(nil, foo{a: 1}))
		assert.Equal(t, foo{a: 1}, ptr.From(&foo{a: 1}))
	})
}

func TestFromMap(t *testing.T) {
	var m map[string]string
	m = ptr.FromMap(m)
	require.NotNil(t, m)
	assert.Equal(t, "", m["foo"])
	m = ptr.FromMap(map[string]string{"foo": "bar"})
	assert.Equal(t, "bar", m["foo"])
}

func TestEqual(t *testing.T) {
	t.Run("nil, nil returns true", func(t *testing.T) {
		assert.True(t, ptr.Equal[int](nil, nil))
	})
	t.Run("nil, not nil and not nil, nil returns false", func(t *testing.T) {
		assert.False(t, ptr.Equal[int](nil, new(1)))
		assert.False(t, ptr.Equal[int](new(1), nil))
	})
	t.Run("&a, &b return *a == *b", func(t *testing.T) {
		assert.False(t, ptr.Equal[int](new(1), new(2)))
		assert.True(t, ptr.Equal[int](new(1), new(1)))
	})
}
