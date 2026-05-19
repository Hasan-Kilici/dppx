package scoring

import (
	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

// Personalization returns a scorer
// using a configurable similarity metric.
func Personalization(
	fn similarity.Func,
) Func {

	if fn == nil {
		fn = similarity.Cosine
	}

	return func(
		query types.Query,
		item types.Item,
	) float64 {

		return fn(
			query.Vector,
			item.Vector,
		)
	}
}