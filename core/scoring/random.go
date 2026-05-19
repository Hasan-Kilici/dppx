package scoring

import (
	"math/rand"

	"github.com/hasan-kilici/dppx/types"
)

// Random produces a random score.
//
// Useful for:
//   - testing
//   - baselines
//   - exploration systems
func Random(
	_ types.Query,
	_ types.Item,
) float64 {

	return rand.Float64()
}