package internal

import "github.com/samber/lo"

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T any, R any](collection []T, iteratee func(item T) R) []R {
	return lo.Map(collection, func(item T, i int) R {
		return iteratee(item)
	})
}
