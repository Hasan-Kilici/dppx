package scoring

import (
	"math/rand"

	"github.com/hasan-kilici/dppx/types"
)

// Random returns an unstructured score for the item.
// This can be used for testing or baseline comparisons.
func Random(item types.Item) float64 {
	return rand.Float64()
}
