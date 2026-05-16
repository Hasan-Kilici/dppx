package similarity

import (
	"math"

	"github.com/hasan-kilici/dppx/types"
)

// Norm computes the Euclidean length of a vector.
func Norm(v types.Vector) float32 {

	var sum float32

	for _, x := range v {
		sum += x * x
	}

	return float32(
		math.Sqrt(float64(sum)),
	)
}
