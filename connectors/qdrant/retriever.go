package qdrant

import (
	"context"
	"fmt"

	qdrant "github.com/qdrant/go-client/qdrant"

	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/core/retriever"
	"github.com/hasan-kilici/dppx/types"
)

var _ retriever.Retriever = (*Retriever)(nil)

type Retriever struct {
	Client     *qdrant.Client
	Collection string
}

func New(
	cfg Config,
) (*Retriever, error) {

	client, err := NewClient(cfg)

	if err != nil {
		return nil, err
	}

	return &Retriever{
		Client:     client,
		Collection: cfg.Collection,
	}, nil
}

func (r *Retriever) Search(
	ctx context.Context,
	query types.Query,
	limit int,
) ([]types.Item, error) {

	limitU64 := uint64(limit)

	resp, err := r.Client.Query(
		ctx,
		&qdrant.QueryPoints{
			CollectionName: r.Collection,

			Query: qdrant.NewQuery(
				query.Vector...,
			),

			Limit: &limitU64,

			WithVectors: &qdrant.WithVectorsSelector{
				SelectorOptions: &qdrant.WithVectorsSelector_Enable{
					Enable: true,
				},
			},

			WithPayload: qdrant.NewWithPayload(
				true,
			),
		},
	)

	if err != nil {
		return nil, err
	}

	items := make([]types.Item, 0, len(resp))

	for _, point := range resp {

		vector := point.Vectors.GetVector().Data

		item := types.Item{
			ID:     fmt.Sprintf("%d", point.Id.GetNum()),
			Vector: vector,
			Norm: similarity.Norm(
				vector,
			),

			Metadata: map[string]any{},
		}

		for k, v := range point.Payload {

			switch val := v.Kind.(type) {

			case *qdrant.Value_StringValue:
				item.Metadata[k] = val.StringValue

			case *qdrant.Value_DoubleValue:
				item.Metadata[k] = val.DoubleValue

			case *qdrant.Value_IntegerValue:
				item.Metadata[k] = val.IntegerValue

			case *qdrant.Value_BoolValue:
				item.Metadata[k] = val.BoolValue
			}
		}

		items = append(items, item)
	}

	return items, nil
}