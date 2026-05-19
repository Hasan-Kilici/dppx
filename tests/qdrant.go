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
    retr, err := qdrant.New(qdrant.Config{
        Host:       "localhost",
        Port:       6334,
        Collection: "blogs",
    })
    if err != nil {
        panic(err)
    }

    eng := engine.New(engine.Config{
        Retriever:     retr,
        CandidatePool: 200,
        Similarity:    similarity.Cosine,
        Scoring: scoring.Combine(
            scoring.Weighted{Func: scoring.Popularity, Weight: 0.6},
            scoring.Weighted{Func: scoring.FreshnessBoost, Weight: 0.3},
            scoring.Weighted{Func: scoring.Personalization, Weight: 0.8},
        ),
        Sampler: sampling.MMR{Lambda: 0.7},
    })

    query := types.Query{UserID: "user-1", Vector: randomVector(768)}
    query.Norm = similarity.Norm(query.Vector)

    results, err := eng.Search(context.Background(), query, 10)
    if err != nil {
        panic(err)
    }

    for _, r := range results {
        fmt.Printf("ITEM=%s SCORE=%.4f\n", r.Item.ID, r.Score)
    }
}

func randomVector(size int) types.Vector {
    v := make(types.Vector, size)
    for i := range v {
        v[i] = 0.5
    }
    return v
}