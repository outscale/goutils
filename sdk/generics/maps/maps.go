package maps

func FromSlice[K comparable, V any, S any](s []S, mapFn func(S) (K, V)) map[K]V {
	m := make(map[K]V, len(s))
	for _, k := range s {
		key, val := mapFn(k)
		m[key] = val
	}
	return m
}
