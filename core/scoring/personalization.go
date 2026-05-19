package scoring

import (
	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

// Personalization measures
// vector similarity between
// query and item embeddings.
func Personalization(
	query types.Query,
	item types.Item,
) float64 {

	return similarity.Cosine(
		query.Vector,
		item.Vector,
		query.Norm,
		item.Norm,
	)
}