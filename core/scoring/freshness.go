package scoring

import (
	"math"
	"time"

	"github.com/hasan-kilici/dppx/types"
)

// FreshnessBoost boosts newer content.
//
// Expected metadata:
//
//	item.Metadata["created_at"] = time.Time
//
// Default half-life is 24 hours.
func FreshnessBoost(
	_ types.Query,
	item types.Item,
) float64 {
	if item.Metadata == nil {
		return 0
	}

	createdAt, ok := item.Metadata["created_at"].(time.Time)

	if !ok {
		return 0
	}

	hours := time.Since(
		createdAt,
	).Hours()

	return math.Exp(
		-hours / 24,
	)
}