package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hasan-kilici/dppx/core/retriever/memory"

	"github.com/hasan-kilici/dppx/core/engine"
	"github.com/hasan-kilici/dppx/core/sampling"
	"github.com/hasan-kilici/dppx/core/scoring"
	"github.com/hasan-kilici/dppx/core/similarity"

	"github.com/hasan-kilici/dppx/types"
)

const (
	UserCount  = 100
	ItemCount  = 100000
	VectorSize = 256
	TopK       = 10
	Topics     = 12
)

var topicCenters []types.Vector

func main() {

	rand.Seed(time.Now().UnixNano())

	topicCenters = generateCenters(
		Topics,
		VectorSize,
	)

	items := generateProducts(
		ItemCount,
	)
	/*
		MEMORY RETRIEVER
	*/
	retr := memory.New(items)
	/*
		DPPX ENGINE
	*/
	eng := engine.New(engine.Config{
		Retriever: retr,

		CandidatePool: 500,

		Similarity: similarity.Cosine,

		Scoring: scoring.Combine(
			scoring.Weighted{
				Func:   scoring.Popularity,
				Weight: 1.0,
			},
			scoring.Weighted{
				Func:   scoring.FreshnessBoost,
				Weight: 0.3,
			},
		),

		Sampler: sampling.MMR{
			Lambda: 0.7,
		},
	})

	var total time.Duration

	for i := 0; i < UserCount; i++ {

		topic := rand.Intn(Topics)

		query := types.Query{
			UserID: fmt.Sprintf(
				"user-%d",
				i,
			),

			Vector: noisy(
				topicCenters[topic],
				0.05,
			),

			Metadata: map[string]any{
				"preferred_topic": topic,
			},
		}

		query.Norm = similarity.Norm(
			query.Vector,
		)

		start := time.Now()

		res, err := eng.Search(
			context.Background(),
			query,
			TopK,
		)

		if err != nil {
			panic(err)
		}

		elapsed := time.Since(start)

		total += elapsed

		fmt.Printf(
			"USER=%d LATENCY=%s RESULTS=%d\n",
			i,
			elapsed,
			len(res),
		)
	}

	fmt.Println()
	fmt.Println("=== BENCHMARK ===")

	fmt.Printf(
		"Avg latency: %s\n",
		total/time.Duration(UserCount),
	)

	fmt.Printf(
		"Total: %s\n",
		total,
	)
}

func generateProducts(
	n int,
) []types.Item {

	items := make([]types.Item, n)

	for i := 0; i < n; i++ {

		topic := rand.Intn(
			Topics,
		)

		vector := noisy(
			topicCenters[topic],
			0.08,
		)

		items[i] = types.Item{
			ID: fmt.Sprintf(
				"item-%d",
				i,
			),

			Vector: vector,

			Norm: similarity.Norm(
				vector,
			),

			Metadata: map[string]any{
				"topic": topic,

				"popularity": pareto(),

				"created_at": time.Now().Add(
					-time.Duration(
						rand.Intn(720),
					) * time.Hour,
				),
			},
		}
	}

	return items
}

func generateCenters(
	k,
	dim int,
) []types.Vector {

	c := make([]types.Vector, k)

	for i := range c {
		c[i] = randomVec(dim)
	}

	return c
}

func randomVec(
	n int,
) types.Vector {

	v := make(types.Vector, n)

	for i := range v {
		v[i] = rand.Float32()
	}

	return v
}

func noisy(
	base types.Vector,
	noise float32,
) types.Vector {

	v := make(
		types.Vector,
		len(base),
	)

	for i := range base {

		v[i] =
			base[i] +
				(rand.Float32()-0.5)*noise
	}

	return v
}

func pareto() float64 {

	alpha := 1.4

	return 1 / math.Pow(
		rand.Float64(),
		1/alpha,
	)
}