package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	qdrantclient "github.com/qdrant/go-client/qdrant"

	"github.com/hasan-kilici/dppx/connectors/qdrant"
	"github.com/hasan-kilici/dppx/core/engine"
	"github.com/hasan-kilici/dppx/core/sampling"
	"github.com/hasan-kilici/dppx/core/scoring"
	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

const (
	CollectionName = "blogs"

	VectorSize   = 768
	InsertCount  = 10000
	CandidatePool = 200
)

func main() {

	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	/*
		Create Qdrant retriever.
	*/

	retr, err := qdrant.New(qdrant.Config{
		Host:       "localhost",
		Port:       6334,
		Collection: CollectionName,
	})

	if err != nil {
		panic(err)
	}

	/*
		Create collection.

		If collection already exists,
		Qdrant will return an error.
	*/

	err = retr.Client.CreateCollection(
		ctx,
		&qdrantclient.CreateCollection{
			CollectionName: CollectionName,

			VectorsConfig: qdrantclient.NewVectorsConfig(
				&qdrantclient.VectorParams{
					Size:     VectorSize,
					Distance: qdrantclient.Distance_Cosine,
				},
			),
		},
	)

	if err != nil {

		fmt.Println(
			"collection may already exist:",
			err,
		)
	}

	/*
		Insert synthetic dataset.

		This simulates a real-world
		recommendation/search corpus.
	*/

	fmt.Println("inserting dataset...")

	points := make(
		[]*qdrantclient.PointStruct,
		0,
		InsertCount,
	)

	for i := 0; i < InsertCount; i++ {

		vector := randomVector(
			VectorSize,
		)

		titleValue, _ := qdrantclient.NewValue(
			fmt.Sprintf(
				"Blog Post %d",
				i,
			),
		)

		categoryValue, _ := qdrantclient.NewValue(
			randomCategory(),
		)

		authorValue, _ := qdrantclient.NewValue(
			fmt.Sprintf(
				"author-%d",
				rand.Intn(100),
			),
		)

		popularityValue, _ := qdrantclient.NewValue(
			rand.Float64(),
		)

		point := &qdrantclient.PointStruct{
			Id: qdrantclient.NewIDNum(
				uint64(i + 1),
			),

			Vectors: qdrantclient.NewVectors(
				vector...,
			),

			Payload: map[string]*qdrantclient.Value{
				"title":      titleValue,
				"category":   categoryValue,
				"author_id":  authorValue,
				"popularity": popularityValue,
			},
		}

		points = append(
			points,
			point,
		)
	}

	_, err = retr.Client.Upsert(
		ctx,
		&qdrantclient.UpsertPoints{
			CollectionName: CollectionName,
			Points:         points,
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("dataset inserted")

	/*
		Create DPPX engine.
	*/

	eng := engine.New(engine.Config{
		Retriever: retr,

		/*
			Candidate pool controls how many
			ANN candidates are retrieved
			before reranking/diversification.
		*/
		CandidatePool: CandidatePool,

		/*
			Vector similarity function.
		*/
		Similarity: similarity.Cosine,

		/*
			Business scoring pipeline.
		*/
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

		/*
			Diversity-aware reranking.
		*/
		Sampler: sampling.MMR{
			Lambda: 0.7,
		},
	})

	/*
		Create query vector.

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
	}

	/*
		Precompute vector norm
		for cosine similarity optimization.
	*/

	query.Norm = similarity.Norm(
		query.Vector,
	)

	/*
		Run retrieval pipeline.
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

	duration := time.Since(start)

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

	/*
		Print benchmark stats.
	*/

	fmt.Println()
	fmt.Println("========== BENCHMARK ==========")

	fmt.Printf(
		"Collection: %s\n",
		CollectionName,
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
	a random embedding vector.

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

/*
	randomCategory generates
	a synthetic content category.
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