package util

func Filter[T any](items []T, filter func(item T) bool) []T {
	result := make([]T, 0)
	for _, v := range items {
		if filter(v) {
			result = append(result, v)
		}
	}
	return result
}

func Delete[S ~[]T, T any](s S, i int) S {
	return append(s[:i], s[i+1:]...)
}
