package sampling

import (
	"math/rand"

	"github.com/hasan-kilici/dppx/types"
)

type Random struct{}

// Sample selects k random items from the scored list.
// The output is useful for baseline comparisons and stochastic selection.
func (r Random) Sample(
	_ types.Query,
	items []types.ScoredItem,
	k int,
) []types.ScoredItem {

	if k >= len(items) {
		return items
	}

	shuffled := make([]types.ScoredItem, len(items))
	copy(shuffled, items)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:k]
}
