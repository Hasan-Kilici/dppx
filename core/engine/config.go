package engine

import (
	"github.com/hasan-kilici/dppx/core/sampling"
	"github.com/hasan-kilici/dppx/core/scoring"
	"github.com/hasan-kilici/dppx/types"
)

// Config defines the search engine behavior through similarity, scoring, and sampling.
type Config struct {
	// Similarity computes how close a query vector is to an item vector.
	Similarity func(
		a types.Vector,
		b types.Vector,
		aNorm float32,
		bNorm float32,
	) float64

	// Scoring is an optional business-level score hook applied alongside similarity.
	Scoring scoring.Func

	// Sampler can post-process scored candidates using a selection strategy.
	Sampler sampling.Sampler

	// CandidatePool is reserved for candidate set sizing or pruning logic.
	CandidatePool int
}
