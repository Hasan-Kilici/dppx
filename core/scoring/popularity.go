package scoring

import "github.com/hasan-kilici/dppx/types"

// Popularity extracts a popularity score from item metadata.
// It returns zero when the popularity field is missing or not numeric.
func Popularity(
	_ types.Query,
	item types.Item,
) float64 {

	if item.Metadata == nil {
		return 0
	}

	popularity, ok := item.Metadata["popularity"].(float64)
	if !ok {
		return 0
	}

	return popularity
}
