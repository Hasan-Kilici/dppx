package main

import (
	"context"
	"fmt"

	"github.com/hasan-kilici/dppx/connectors/qdrant"
	"github.com/hasan-kilici/dppx/core/engine"
	"github.com/hasan-kilici/dppx/core/sampling"
	"github.com/hasan-kilici/dppx/core/scoring"
	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

func main() {

	/*
		Create a Qdrant retriever.

		This connector is responsible for:
		- communicating with Qdrant
		- retrieving ANN candidates
		- converting Qdrant points into DPPX items
	*/

	retr, err := qdrant.New(qdrant.Config{
		Host:       "localhost",
		Port:       6334,
		Collection: "blogs",
	})

	if err != nil {
		panic(err)
	}

	/*
		Create the DPPX engine.

		The engine handles:
		- similarity scoring
		- business scoring
		- reranking
		- diversity sampling
	*/

	eng := engine.New(engine.Config{
		Retriever: retr,

		// Candidate pool size retrieved from ANN layer
		CandidatePool: 200,

		// Vector similarity function
		Similarity: similarity.Cosine,

		// Business scoring pipeline
		Scoring: scoring.Combine(
			scoring.Weighted{
				Func:   scoring.Popularity,
				Weight: 0.6,
			},
			scoring.Weighted{
				Func:   scoring.FreshnessBoost,
				Weight: 0.3,
			},
			scoring.Weighted{
				Func:   scoring.Personalization,
				Weight: 0.8,
			},
		),

		// Diversity-aware reranking
		Sampler: sampling.MMR{
			Lambda: 0.7,
		},
	})

	/*
		Create a query vector.

		In real-world systems this vector usually comes from:
		- OpenAI embeddings
		- sentence-transformers
		- BGE
		- E5
		- Cohere embeddings
		- custom embedding models
	*/

	query := types.Query{
		UserID: "user-1",

		Vector: randomVector(768),
	}

	// Precomputed vector norm for faster cosine similarity
	query.Norm = similarity.Norm(
		query.Vector,
	)

	/*
		Run retrieval + reranking pipeline.
	*/

	results, err := eng.Search(
		context.Background(),
		query,
		10,
	)

	if err != nil {
		panic(err)
	}

	/*
		Print final results.
	*/

	for _, result := range results {

		fmt.Printf(
			"ITEM=%s SCORE=%.4f CATEGORY=%v\n",
			result.Item.ID,
			result.Score,
			result.Item.Metadata["category"],
		)
	}
}

/*
	Example helper.

	Replace this with real embeddings in production.
*/

func randomVector(
	size int,
) types.Vector {

	v := make(
		types.Vector,
		size,
	)

	for i := range v {
		v[i] = 0.5
	}

	return v
}