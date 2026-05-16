package scoring

import (
	"math"

	"github.com/hasan-kilici/dppx/types"
)

// Norm computes the Euclidean norm of an item's vector.
func Norm(
	query types.Query,
	item types.Item,
) float64 {

	var sum float32

	for _, v := range item.Vector {
		sum += v * v
	}

	return math.Sqrt(float64(sum))
}
