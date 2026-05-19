package similarity

import "github.com/hasan-kilici/dppx/types"

// Dot computes dot-product similarity.
func Dot(
	a types.Vector,
	b types.Vector,
) float64 {

	if len(a) != len(b) {
		return 0
	}

	var dot float32

	for i := range a {
		dot += a[i] * b[i]
	}

	return float64(dot)
}