package similarity

import (
	"math"

	"github.com/hasan-kilici/dppx/types"
)

// Cosine computes cosine similarity
// between two dense vectors.
func Cosine(
	a types.Vector,
	b types.Vector,
) float64 {

	if len(a) != len(b) {
		return 0
	}

	var (
		dot   float32
		normA float32
		normB float32
	)

	for i := range a {

		dot += a[i] * b[i]

		normA += a[i] * a[i]

		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float64(
		dot / float32(
			math.Sqrt(float64(normA*normB)),
		),
	)
}