/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package ptr

// To returns a pointer to a value.
// Deprecated: use new.
//
//go:fix inline
func To[T any](t T) *T {
	return new(t)
}

// From returns def (if specified) or a zero value if nil or the value referenced by the pointer.
func From[T any](t *T, def ...T) T {
	switch {
	case t != nil:
		return *t
	case len(def) > 0:
		return def[0]
	default:
		var tt T
		return tt
	}
}

// From returns an empty map if nil or the map value.
func FromMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return make(map[K]V)
	}
	return m
}

// Equal returns true if both pointers are nil or point to the same value.
func Equal[T comparable](a, b *T) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a != nil && b != nil:
		return *a == *b
	default:
		return false
	}
}
