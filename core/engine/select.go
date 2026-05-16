package engine

import "github.com/hasan-kilici/dppx/types"

// Select returns the top-k items only, omitting score metadata.
func (e *Engine) Select(
	items []types.Item,
	k int,
) []types.Item {

	results := e.Search(
		types.Query{},
		items,
		k,
	)

	final := make([]types.Item, 0, len(results))

	for _, result := range results {
		final = append(final, result.Item)
	}

	return final
}
