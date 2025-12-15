/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package batch

func sliceToMap[T any, K comparable, V any](in []T, fn func(t T) (K, V)) map[K]V {
	out := make(map[K]V, len(in))
	for i := range in {
		k, v := fn(in[i])
		out[k] = v
	}
	return out
}
