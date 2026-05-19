package sampling

import (
	"math"

	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

// similarityPair caches pairwise similarities.
type similarityPair struct {
	A string
	B string
}

// MMR implements Maximal Marginal Relevance.
type MMR struct {

	// Lambda balances:
	// relevance vs diversity
	Lambda float64

	// Similarity metric used
	// for diversity comparison.
	Similarity similarity.Func
}

func getSimilarity(
	cache map[similarityPair]float64,

	fn similarity.Func,

	a types.ScoredItem,
	b types.ScoredItem,
) float64 {

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

	sim := fn(
		a.Item.Vector,
		b.Item.Vector,
	)

	cache[key] = sim

	return sim
}

// Sample selects diverse results.
func (m MMR) Sample(
	_ types.Query,
	items []types.ScoredItem,
	k int,
) []types.ScoredItem {

	if k >= len(items) {
		return items
	}

	/*Default similarity*/
	if m.Similarity == nil {
		m.Similarity = similarity.Cosine
	}

	cache := make(
		map[similarityPair]float64,
	)

	selected := make(
		[]types.ScoredItem,
		0,
		k,
	)

	remaining := make(
		[]types.ScoredItem,
		len(items),
	)

	copy(remaining, items)

	selected = append(
		selected,
		remaining[0],
	)

	remaining = remaining[1:]

	for len(selected) < k &&
		len(remaining) > 0 {

		bestIdx := 0

		bestScore := math.Inf(-1)

		for i, candidate := range remaining {
			maxSimilarity := 0.0
			for _, chosen := range selected {

				sim := getSimilarity(
					cache,
					m.Similarity,
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