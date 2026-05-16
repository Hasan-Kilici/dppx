package sampling

import "github.com/hasan-kilici/dppx/types"

// Sampler selects a subset of scored items using a custom strategy.
type Sampler interface {
	Sample(
		query types.Query,
		items []types.ScoredItem,
		k int,
	) []types.ScoredItem
}
