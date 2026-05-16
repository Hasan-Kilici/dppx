package scoring

import (
	"time"

	"github.com/hasan-kilici/dppx/types"
)

// FreshnessBoost gives newer items a higher score with a 24-hour half-life decay.
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

	hours := time.Since(createdAt).Hours()

	// newer content => higher score
	return 1 / (1 + hours/24)
}
