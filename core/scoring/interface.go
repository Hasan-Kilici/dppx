package scoring

import "github.com/hasan-kilici/dppx/types"

// Func defines a scoring callback to compute an item score for a given query.
type Func func(query types.Query, item types.Item) float64
