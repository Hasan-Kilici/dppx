package retriever

import (
	"context"

	"github.com/hasan-kilici/dppx/types"
)

type Memory struct {
	Items []types.Item
}

func (m Memory) Retrieve(
	_ context.Context,
	_ types.Query,
	k int,
) ([]types.Item, error) {

	if k > len(m.Items) {
		k = len(m.Items)
	}

	return m.Items[:k], nil
}