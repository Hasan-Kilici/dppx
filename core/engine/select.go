package engine

import (
	"context"

	"github.com/hasan-kilici/dppx/types"
)

// Select returns only items without scores.
func (e *Engine) Select(
	ctx context.Context,
	query types.Query,
	k int,
) ([]types.Item, error) {

	results, err := e.Search(
		ctx,
		query,
		k,
	)

	if err != nil {
		return nil, err
	}

	final := make([]types.Item, 0, len(results))

	for _, result := range results {
		final = append(final, result.Item)
	}

	return final, nil
}