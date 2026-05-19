package engine

import (
	"github.com/hasan-kilici/dppx/core/retriever"
	"github.com/hasan-kilici/dppx/core/sampling"
	"github.com/hasan-kilici/dppx/core/scoring"
	"github.com/hasan-kilici/dppx/core/similarity"
)

// Config defines engine behavior.
type Config struct {

	// Retriever fetches ANN candidates from vector DBs.
	Retriever retriever.Retriever

	// Similarity computes vector similarity.
	Similarity similarity.Func

	// Optional business scoring layer.
	Scoring scoring.Func

	// Diversity / reranking stage.
	Sampler sampling.Sampler

	// Number of retrieved candidates before reranking.
	CandidatePool int
}