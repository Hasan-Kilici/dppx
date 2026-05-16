package rerank

import "github.com/hasan-kilici/dppx/types"

// Reranker defines a secondary ranking stage that can reorder or prune items.
type Reranker interface {
	Rerank(items []types.Item, k int) []types.Item
}
