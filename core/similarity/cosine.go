package similarity

import "github.com/hasan-kilici/dppx/types"

// Cosine computes cosine similarity between two pre-normalized vectors.
// The caller must provide precomputed norms to avoid repeated expensive sqrt operations.
func Cosine(
	a types.Vector,
	b types.Vector,
	aNorm float32,
	bNorm float32,
) float64 {

	if len(a) != len(b) {
		return 0
	}

	if aNorm == 0 || bNorm == 0 {
		return 0
	}

	var dot float32

	for i := range a {
		dot += a[i] * b[i]
	}

	return float64(
		dot / (aNorm * bNorm),
	)
}
