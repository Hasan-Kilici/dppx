package retriever

import (
	"context"

	"github.com/hasan-kilici/dppx/types"
)

type Retriever interface {
	Search(
		ctx context.Context,
		query types.Query,
		limit int,
	) ([]types.Item, error)
}