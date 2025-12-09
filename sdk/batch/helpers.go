package batch

func sliceToMap[T any, K comparable, V any](in []T, fn func(t T) (K, V)) map[K]V {
	out := make(map[K]V, len(in))
	for i := range in {
		k, v := fn(in[i])
		out[k] = v
	}
	return out
}
