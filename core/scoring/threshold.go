package scoring

import "github.com/hasan-kilici/dppx/types"

// Threshold applies a minimum score gate.
//
// If the base scorer returns a value
// below threshold, zero is returned.
func Threshold(
	base Func,
	threshold float64,
) Func {

	return func(
		query types.Query,
		item types.Item,
	) float64 {

		if base == nil {
			return 0
		}

		score := base(
			query,
			item,
		)

		if score < threshold {
			return 0
		}

		return score
	}
}