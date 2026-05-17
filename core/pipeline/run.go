package pipeline

import (
	"context"
	"sort"

	"github.com/hasan-kilici/dppx/fusion"
	"github.com/hasan-kilici/dppx/types"
)

func (p *Pipeline) Run(
	ctx context.Context,
	query types.Query,
	k int,
) ([]types.ScoredItem, error) {

	candidates, err := p.Retriever.Retrieve(
		ctx,
		query,
		1000,
	)

	if err != nil {
		return nil, err
	}

	scored := make([]types.ScoredItem, 0, len(candidates))

	for _, item := range candidates {

		similarityScore := p.Similarity(
			query.Vector,
			item.Vector,
			query.Norm,
			item.Norm,
		)

		customScore := float32(0)

		if p.Scoring != nil {
			customScore = p.Scoring(
				query,
				item,
			)
		}

		finalScore := p.Fusion.Combine(
			fusion.Scores{
				Similarity: similarityScore,
				Custom: customScore,
			},
		)

		scored = append(scored, types.ScoredItem{
			Item: item,
			Score: finalScore,
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if p.Sampler != nil {
		scored = p.Sampler.Sample(
			query,
			scored,
			k,
		)
	}

	if k > len(scored) {
		k = len(scored)
	}

	return scored[:k], nil
}