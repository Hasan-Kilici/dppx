package sampling

import "github.com/hasan-kilici/dppx/types"

// TopK returns the first k items without additional reordering.
// Use this when the input list is already scored and sorted.
type TopK struct{}

func (t TopK) Sample(
	_ types.Query,
	items []types.ScoredItem,
	k int,
) []types.ScoredItem {

	if k > len(items) {
		k = len(items)
	}

	return items[:k]
}
