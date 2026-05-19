package qdrant

import (
	qdrant "github.com/qdrant/go-client/qdrant"
)

func Match(
	key string,
	value any,
) *qdrant.Condition {
	switch v := value.(type) {
	case string:
		return qdrant.NewMatch(
			key,
			v,
		)
	case int:
		return qdrant.NewMatchInt(
			key,
			int64(v),
		)
	case int64:
		return qdrant.NewMatchInt(
			key,
			v,
		)
	case bool:
		return qdrant.NewMatchBool(
			key,
			v,
		)
	default:
		panic("unsupported qdrant match type")
	}
}

func FilterMust(
	conditions ...*qdrant.Condition,
) *qdrant.Filter {

	return &qdrant.Filter{
		Must: conditions,
	}
}