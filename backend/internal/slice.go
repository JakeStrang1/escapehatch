package internal

import "github.com/samber/lo"

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T any, R any](collection []T, iteratee func(item T) R) []R {
	return lo.Map(collection, func(item T, i int) R {
		return iteratee(item)
	})
}

// Filter manipulates a slice by removing any elements where iteratee(element) == false
func Filter[T any](collection []T, iteratee func(item T) bool) []T {
	return lo.FilterMap(collection, func(item T, _ int) (T, bool) {
		return item, iteratee(item)
	})
}
