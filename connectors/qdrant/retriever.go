package qdrant

import (
	"context"

	qdrant "github.com/qdrant/go-client/qdrant"

	"github.com/hasan-kilici/dppx/core/retriever"
	"github.com/hasan-kilici/dppx/types"
)

var _ retriever.Retriever = (*Retriever)(nil)

type Retriever struct {
	Client     *Client
	Collection string
}

func New(
	cfg Config,
) (*Retriever, error) {

	client, err := NewClient(
		cfg,
	)

	if err != nil {
		return nil, err
	}

	return &Retriever{
		Client: client,
		Collection: cfg.Collection,
	}, nil
}

/*
	Default DPPX retrieval path.

	This is intentionally minimal and generic.

	Advanced Qdrant functionality should use:
	- SearchWithOptions
	- Raw client access
*/
func (r *Retriever) Search(
	ctx context.Context,
	query types.Query,
	limit int,
) ([]types.Item, error) {
	opts := DefaultSearchOptions()
	opts.Limit = limit

	return r.SearchWithOptions(
		ctx,
		query,
		opts,
	)
}

/*Production-grade configurable search.*/
func (r *Retriever) SearchWithOptions(
	ctx context.Context,
	query types.Query,
	opts SearchOptions,
) ([]types.Item, error) {
	limit := uint64(
		opts.Limit,
	)

	req := &qdrant.QueryPoints{
		CollectionName: r.Collection,

		Query: qdrant.NewQuery(
			query.Vector...,
		),
		Limit: &limit,
	}

	/*Optional payload retrieval*/
	if opts.WithPayload {
		req.WithPayload =
			qdrant.NewWithPayload(
				true,
			)
	}

	/*Optional vector retrieval*/
	if opts.WithVectors {
		req.WithVectors =
			qdrant.NewWithVectors(
				true,
			)
	}

	/*Optional metadata filters*/
	if opts.Filter != nil {
		req.Filter = opts.Filter
	}

	/*Optional ANN tuning params*/
	if opts.Params != nil {
		req.Params = opts.Params
	}

	/*Optional score threshold*/

	if opts.ScoreThreshold != nil {
		req.ScoreThreshold =
			opts.ScoreThreshold
	}

	resp, err := r.Client.Raw().Query(
		ctx,
		req,
	)

	if err != nil {
		return nil, err
	}

	items := make(
		[]types.Item,
		0,
		len(resp),
	)

	for _, point := range resp {
		item := mapPoint(
			point,
		)

		items = append(
			items,
			item,
		)
	}

	return items, nil
}