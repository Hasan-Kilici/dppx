package pool

import (
	"sync"

	"github.com/hasan-kilici/dppx/types"
)

// ScoredItemPool reuses ScoredItem objects to reduce allocation pressure in hot workloads.
var ScoredItemPool = sync.Pool{
	New: func() any {
		return new(types.ScoredItem)
	},
}
