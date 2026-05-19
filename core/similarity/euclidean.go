package similarity

import (
	"math"

	"github.com/hasan-kilici/dppx/types"
)

// Euclidean computes euclidean similarity.
// Higher value means more similar.
func Euclidean(
	a types.Vector,
	b types.Vector,
) float64 {

	if len(a) != len(b) {
		return 0
	}

	var sum float64

	for i := range a {

		diff := float64(a[i] - b[i])

		sum += diff * diff
	}

	distance := math.Sqrt(sum)

	// convert distance -> similarity
	return 1 / (1 + distance)
}