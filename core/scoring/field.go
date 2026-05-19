package scoring

import "github.com/hasan-kilici/dppx/types"

// Field returns a scoring function
// that extracts a numeric value
// from metadata dynamically.
//
// Example:
//
//	scoring.Field("ctr")
//	scoring.Field("quality")
//	scoring.Field("watch_time")
func Field(
	key string,
) Func {

	return func(
		_ types.Query,
		item types.Item,
	) float64 {

		if item.Metadata == nil {
			return 0
		}

		value, ok := item.Metadata[key]

		if !ok {
			return 0
		}

		switch v := value.(type) {

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
}