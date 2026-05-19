package scoring

import "github.com/hasan-kilici/dppx/types"

// Popularity extracts a popularity score
// from item metadata.
//
// Expected metadata:
//
//	item.Metadata["popularity"] = float64
func Popularity(
	_ types.Query,
	item types.Item,
) float64 {
	if item.Metadata == nil {
		return 0
	}

	switch v := item.Metadata["popularity"].(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}

	return 0
}