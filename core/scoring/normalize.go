package scoring

import "github.com/hasan-kilici/dppx/types"

// Normalize clamps scorer output
// into the range [0, 1].
func Normalize(
	base Func,
	max float64,
) Func {

	return func(
		query types.Query,
		item types.Item,
	) float64 {

		if base == nil || max <= 0 {
			return 0
		}

		score := base(
			query,
			item,
		)

		if score <= 0 {
			return 0
		}

		normalized := score / max

		if normalized > 1 {
			return 1
		}

		return normalized
	}
}