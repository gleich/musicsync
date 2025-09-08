package utils

// Batch splits items into consecutive groups of at most size elements,
// preserving the original order. The final group may contain fewer than
// size elements if len(items) is not a multiple of size.
//
// If size <= 0, Batch returns nil.
//
// Note: if items is empty, Batch returns a slice containing a single
// empty batch (i.e., [][]T{nil}).
func Batch[T any](items []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	groups := make([][]T, 1)
	for _, item := range items {
		if len(groups[len(groups)-1]) >= size {
			groups = append(groups, []T{})
		}
		groups[len(groups)-1] = append(groups[len(groups)-1], item)
	}
	return groups
}
