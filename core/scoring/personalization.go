package scoring

import (
	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

// Personalization scores an item by measuring similarity to the query vector.
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
