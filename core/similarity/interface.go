package similarity

import "github.com/hasan-kilici/dppx/types"

// Func defines a vector similarity metric.
type Func func(
	query,
	item types.Vector,
) float64