package memory

import (
	"context"

	"github.com/hasan-kilici/dppx/core/retriever"
	"github.com/hasan-kilici/dppx/types"
)

var _ retriever.Retriever = (*Retriever)(nil)

type Retriever struct {
	Items []types.Item
}

func New(
	items []types.Item,
) *Retriever {

	return &Retriever{
		Items: items,
	}
}

func (r *Retriever) Search(
	_ context.Context,
	_ types.Query,
	limit int,
) ([]types.Item, error) {

	if limit > len(r.Items) {
		limit = len(r.Items)
	}

	return r.Items[:limit], nil
}