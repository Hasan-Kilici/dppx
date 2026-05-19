package scoring

import "github.com/hasan-kilici/dppx/types"

// Conditional applies a scorer
// only when predicate returns true.
func Conditional(
	predicate func(
		query types.Query,
		item types.Item,
	) bool,

	scorer Func,
) Func {

	return func(
		query types.Query,
		item types.Item,
	) float64 {

		if predicate == nil || scorer == nil {
			return 0
		}

		if !predicate(
			query,
			item,
		) {
			return 0
		}

		return scorer(
			query,
			item,
		)
	}
}