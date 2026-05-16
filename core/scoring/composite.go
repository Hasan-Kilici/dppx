package scoring

import "github.com/hasan-kilici/dppx/types"

// Weighted associates a scoring function with a relative weight.
// Use this to create composite scoring logic from multiple signals.
type Weighted struct {
	Func   Func
	Weight float64
}

// Combine merges multiple weighted scorers into a single scoring function.
func Combine(
	scorers ...Weighted,
) Func {

	return func(
		query types.Query,
		item types.Item,
	) float64 {

		var total float64

		for _, scorer := range scorers {
			total += scorer.Func(query, item) * scorer.Weight
		}

		return total
	}
}
