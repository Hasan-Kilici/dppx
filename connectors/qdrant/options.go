package qdrant

import qdrant "github.com/qdrant/go-client/qdrant"

type SearchOptions struct {
	Limit int

	WithPayload bool
	WithVectors bool

	Filter *qdrant.Filter

	Params *qdrant.SearchParams

	ScoreThreshold *float32
}

func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		Limit: 10,

		WithPayload: true,
		WithVectors: true,
	}
}