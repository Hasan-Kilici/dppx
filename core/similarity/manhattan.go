package similarity

import (
	"math"

	"github.com/hasan-kilici/dppx/types"
)

// Manhattan computes L1 similarity.
//
// Higher value means more similar.
func Manhattan(
	a types.Vector,
	b types.Vector,
) float64 {

	if len(a) != len(b) {
		return 0
	}

	var distance float64

	for i := range a {

		distance += math.Abs(
			float64(a[i] - b[i]),
		)
	}

	// convert distance -> similarity
	return 1 / (1 + distance)
}