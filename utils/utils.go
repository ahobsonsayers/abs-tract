package utils

import (
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
)

func ToPointer[T any](value T) *T {
	return &value
}

// Intersection returns the elements in slice 1 that also
// appear in slice 2. Returned slice order will match that
// of slice 1. If slice 1 has duplicates of a value that is
// present in slice 2, all duplicates will be included.
func Intersection[T comparable](slice1, slice2 []T) []T {
	set2 := mapset.NewSet(slice2...)

	intersection := make([]T, 0, len(slice1))
	for _, elem := range slice1 {
		if set2.Contains(elem) {
			intersection = append(intersection, elem)
		}
	}

	return slices.Clip(intersection)
}
