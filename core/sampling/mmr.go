package sampling

import (
	"math"

	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

// similarityPair is used as a cache key for symmetric pairwise similarity lookups.
type similarityPair struct {
	A string
	B string
}

// MMR implements Maximal Marginal Relevance sampling to balance relevance and diversity.
type MMR struct {
	Lambda float64
}

func getSimilarity(
	cache map[similarityPair]float64,
	a types.ScoredItem,
	b types.ScoredItem,
) float64 {
	// Reuse symmetric similarity values to avoid redundant cosine computations.

	key := similarityPair{
		A: a.Item.ID,
		B: b.Item.ID,
	}

	reverse := similarityPair{
		A: b.Item.ID,
		B: a.Item.ID,
	}

	if sim, ok := cache[key]; ok {
		return sim
	}

	if sim, ok := cache[reverse]; ok {
		return sim
	}

	sim := similarity.Cosine(
		a.Item.Vector,
		b.Item.Vector,
		a.Item.Norm,
		b.Item.Norm,
	)

	cache[key] = sim

	return sim
}

// Sample chooses a diverse subset of scored items using the MMR objective.
// Higher lambda favors relevance, lower lambda favors diversity.
func (m MMR) Sample(
	_ types.Query,
	items []types.ScoredItem,
	k int,
) []types.ScoredItem {

	if k >= len(items) {
		return items
	}

	cache := make(map[similarityPair]float64)

	selected := make([]types.ScoredItem, 0, k)

	remaining := make([]types.ScoredItem, len(items))
	copy(remaining, items)

	selected = append(selected, remaining[0])
	remaining = remaining[1:]

	for len(selected) < k && len(remaining) > 0 {

		bestIdx := 0
		bestScore := math.Inf(-1)

		for i, candidate := range remaining {

			maxSimilarity := 0.0

			for _, chosen := range selected {

				sim := getSimilarity(
					cache,
					candidate,
					chosen,
				)

				if sim > maxSimilarity {
					maxSimilarity = sim
				}
			}

			mmrScore :=
				(m.Lambda * candidate.Score) -
					((1 - m.Lambda) * maxSimilarity)

			if mmrScore > bestScore {
				bestScore = mmrScore
				bestIdx = i
			}
		}

		selected = append(
			selected,
			remaining[bestIdx],
		)

		remaining = append(
			remaining[:bestIdx],
			remaining[bestIdx+1:]...,
		)
	}

	return selected
}
