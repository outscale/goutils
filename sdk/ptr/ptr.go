/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package ptr

func To[T any](t T) *T {
	return &t
}

func From[T any](t *T) T {
	if t == nil {
		var tt T
		return tt
	}
	return *t
}
