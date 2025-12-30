/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package ptr

// To returns a pointer to a value.
func To[T any](t T) *T {
	return &t
}

// From returns a zero value if nil or the value referenced by the pointer.
func From[T any](t *T) T {
	if t == nil {
		var tt T
		return tt
	}
	return *t
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
