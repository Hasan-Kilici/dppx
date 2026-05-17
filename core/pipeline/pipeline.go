package pipeline

import (
	"context"

	"github.com/hasan-kilici/dppx/fusion"
	"github.com/hasan-kilici/dppx/retriever"

	"github.com/hasan-kilici/dppx/core/sampling"
	"github.com/hasan-kilici/dppx/core/scoring"
	"github.com/hasan-kilici/dppx/core/similarity"

	"github.com/hasan-kilici/dppx/types"
)

type Pipeline struct {
	Retriever retriever.Retriever
	Similarity similarity.Func
	Scoring scoring.Func
	Fusion fusion.Fusion
	Sampler sampling.Sampler
}