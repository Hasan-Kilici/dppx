package scoring

import "github.com/hasan-kilici/dppx/types"

// Boost returns a scorer that
// multiplies another scorer
// by a constant value.
func Boost(
	base Func,
	multiplier float64,
) Func {

	return func(
		query types.Query,
		item types.Item,
	) float64 {

		if base == nil {
			return 0
		}

		return base(
			query,
			item,
		) * multiplier
	}
}