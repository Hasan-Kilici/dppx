package scoring

import "github.com/hasan-kilici/dppx/types"

// Func defines a scoring callback.
//
// The returned value is expected to be normalized
// between 0.0 and 1.0 when possible.
type Func func(
	query types.Query,
	item types.Item,
) float64