package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hasan-kilici/dppx/connectors/mongodb"

	"github.com/hasan-kilici/dppx/core/engine"
	"github.com/hasan-kilici/dppx/core/sampling"
	"github.com/hasan-kilici/dppx/core/scoring"
	"github.com/hasan-kilici/dppx/core/similarity"

	"github.com/hasan-kilici/dppx/types"
)

const (
	VectorSize = 768
)

func main() {

	ctx := context.Background()

	rand.Seed(
		time.Now().UnixNano(),
	)

	/*
		Create MongoDB retriever.

		This example uses MongoDB Atlas
		Vector Search integration.

		Collection schema example:

		{
			"_id": "...",
			"title": "...",
			"category": "...",
			"embedding": [...]
		}
	*/

	retr, err := mongodb.New(
		ctx,
		mongodb.Config{
			URI: "mongodb://localhost:27017",

			Database: "blogs",

			Collection: "blogs",

			/*
				Field containing embeddings.
			*/
			VectorField: "embedding",

			/*
				Atlas Vector Search index name.
			*/
			IndexName: "vector_index",

			/*
				false:
				use local retrieval mode

				true:
				use Atlas Vector Search
			*/
			UseAtlasVectorSearch: false,
		},
	)

	if err != nil {
		panic(err)
	}

	/*
		Create DPPX engine.
	*/

	eng := engine.New(
		engine.Config{
			Retriever: retr,

			/*
				Number of retrieved
				candidates before reranking.
			*/
			CandidatePool: 200,

			/*
				Main vector similarity metric.
			*/
			Similarity: similarity.Cosine,

			/*
				Business scoring pipeline.
			*/
			Scoring: scoring.Combine(
				scoring.Weighted{
					Func: scoring.Popularity,
					Weight: 0.6,
				},

				scoring.Weighted{
					Func: scoring.FreshnessBoost,
					Weight: 0.3,
				},

				scoring.Weighted{
					Func: scoring.Personalization(
						similarity.Cosine,
					),
					Weight: 0.8,
				},
			),

			/*
				Diversity-aware reranking.
			*/
			Sampler: sampling.MMR{
				Lambda: 0.7,

				Similarity:
					similarity.Cosine,
			},
		},
	)

	/*
		Query vector.

		In production this usually comes from:
		- OpenAI embeddings
		- BGE
		- E5
		- sentence-transformers
		- Cohere embeddings
	*/

	query := types.Query{
		UserID: "user-1",

		Vector: randomVector(
			VectorSize,
		),

		Metadata: map[string]any{
			"preferred_category": "ai",
		},
	}

	/*
		Run retrieval pipeline.
	*/

	results, err := eng.Search(
		ctx,
		query,
		10,
	)

	if err != nil {
		panic(err)
	}

	/*
		Print results.
	*/

	fmt.Println()
	fmt.Println("========== RESULTS ==========")

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
	randomVector creates
	a synthetic embedding vector.

	Replace this with real embeddings
	in production systems.
*/

func randomVector(
	size int,
) types.Vector {

	v := make(
		types.Vector,
		size,
	)

	for i := range v {
		v[i] = rand.Float32()
	}

	return v
}