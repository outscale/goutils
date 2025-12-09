package maps

func Map[K comparable, V any, KK comparable, VV any](in map[K]V, fn func(k K, v V) (KK, VV, bool)) map[KK]VV {
	out := make(map[KK]VV, len(in))
	for k, v := range in {
		nk, nv, ok := fn(k, v)
		if !ok {
			continue
		}
		out[nk] = nv
	}
	return out
}

func FromSlice[K comparable, V any, S any](s []S, fn func(S) (K, V)) map[K]V {
	m := make(map[K]V, len(s))
	for _, k := range s {
		key, val := fn(k)
		m[key] = val
	}
	return m
}
