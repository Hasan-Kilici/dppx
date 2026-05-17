package main

import (
	"context"
	"fmt"

	"github.com/hasan-kilici/dppx/core/engine"
	mem "github.com/hasan-kilici/dppx/core/retriever/memory"
	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

func main() {

	items := []types.Item{
		{
			ID: "a",
			Vector: types.Vector{
				0.1,
				0.2,
				0.3,
			},
			Norm: 0.374,
		},
		{
			ID: "b",
			Vector: types.Vector{
				0.2,
				0.1,
				0.4,
			},
			Norm: 0.469,
		},
	}

	retr := mem.New(items)

	cfg := engine.Config{
		Retriever: retr,

		CandidatePool: 10,

		Similarity: similarity.Cosine,
	}

	eng := engine.New(cfg)

	query := types.Query{
		Vector: types.Vector{
			0.1,
			0.2,
			0.3,
		},

		Norm: 0.374,
	}

	results, err := eng.Search(
		context.Background(),
		query,
		10,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(results)
}