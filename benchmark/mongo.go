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

	"go.mongodb.org/mongo-driver/bson"
)

const (
	VectorSize    = 768
	InsertCount   = 10000
	CandidatePool = 200
)

func main() {

	ctx := context.Background()

	rand.Seed(
		time.Now().UnixNano(),
	)

	/*
		MongoDB retriever
	*/

	retr, err := mongodb.New(
		ctx,
		mongodb.Config{
			URI: "mongodb://localhost:27017",
			Database: "blogs",
			Collection: "blogs",
			VectorField: "embedding",

			/*
				Atlas Vector Search index name.

				Used only if Atlas vector
				search is enabled.
			*/
			IndexName: "vector_index",
		},
	)

	if err != nil {
		panic(err)
	}

	/*
		Raw Mongo collection access.

		Useful for:
		- inserting documents
		- custom aggregation pipelines
		- admin operations
	*/

	coll := retr.Client().
				Raw().
				Database("blogs").
				Collection("blogs")

	/*
		Clean old dataset
	*/

	_, err = coll.DeleteMany(
		ctx,
		bson.M{},
	)

	if err != nil {
		panic(err)
	}

	/*
		Insert synthetic dataset
	*/

	fmt.Println("inserting dataset...")

	docs := make(
		[]any,
		0,
		InsertCount,
	)

	for i := 0; i < InsertCount; i++ {

		vector := randomVector(
			VectorSize,
		)

		doc := bson.M{
			"title": fmt.Sprintf(
				"Blog Post %d",
				i,
			),

			"category": randomCategory(),

			"author_id": fmt.Sprintf(
				"author-%d",
				rand.Intn(100),
			),

			"popularity": rand.Float64(),

			"created_at": time.Now().Add(
				-time.Duration(
					rand.Intn(720),
				) * time.Hour,
			),

			/*
				Vector embedding field.

				This field name must match:
				Config.VectorField
			*/
			"embedding": vector,
		}

		docs = append(
			docs,
			doc,
		)
	}

	_, err = coll.InsertMany(
		ctx,
		docs,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("dataset inserted")

	/*
		DPPX engine
	*/

	eng := engine.New(
		engine.Config{
			Retriever: retr,

			CandidatePool:
				CandidatePool,

			Similarity:
				similarity.Cosine,

			/*
				Business scoring layer
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
				Diversity-aware reranking
			*/
			Sampler: sampling.MMR{
				Lambda: 0.7,

				Similarity:
					similarity.Cosine,
			},
		},
	)

	/*
		Query vector
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
		Run benchmark
	*/

	fmt.Println()
	fmt.Println("running benchmark...")

	start := time.Now()

	results, err := eng.Search(
		ctx,
		query,
		10,
	)

	if err != nil {
		panic(err)
	}

	duration := time.Since(
		start,
	)

	/*
		Print results
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

	/*
		Print benchmark stats
	*/

	fmt.Println()
	fmt.Println("========== BENCHMARK ==========")

	fmt.Printf(
		"Database: %s\n",
		"blogs",
	)

	fmt.Printf(
		"Collection: %s\n",
		"blogs",
	)

	fmt.Printf(
		"Inserted Items: %d\n",
		InsertCount,
	)

	fmt.Printf(
		"Vector Size: %d\n",
		VectorSize,
	)

	fmt.Printf(
		"Candidate Pool: %d\n",
		CandidatePool,
	)

	fmt.Printf(
		"Latency: %s\n",
		duration,
	)
}

/*
	randomVector creates
	a synthetic embedding vector.
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

/*
	randomCategory generates
	a synthetic category.
*/

func randomCategory() string {

	categories := []string{
		"golang",
		"ai",
		"backend",
		"distributed-systems",
		"devops",
		"databases",
		"vector-search",
		"recommendation",
	}

	return categories[
		rand.Intn(
			len(categories),
		),
	]
}