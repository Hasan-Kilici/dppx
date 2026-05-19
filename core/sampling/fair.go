package sampling

import (
	"sort"

	"github.com/hasan-kilici/dppx/types"
)

/*
	Fair sampler performs diversity-aware
	group balancing.

	Example use-cases:
	- category balancing
	- author balancing
	- language balancing
	- marketplace seller balancing
	- political/source balancing
	- feed deduplication

	The sampler tries to prevent
	a single group from dominating
	the final ranking.

	Example:
	Without Fair:
		AI AI AI AI Backend AI AI

	With Fair(group=category):
		AI Backend DevOps AI Database

	Group values are extracted from:
		item.Metadata[GroupKey]
*/

type Fair struct {

	/*
		GroupKey defines which metadata
		field will be used for balancing.

		Examples:
		- "category"
		- "author_id"
		- "language"
		- "seller_id"
	*/
	GroupKey string

	/*
		MaxPerGroup limits how many items
		can appear from the same group.

		Example:
		MaxPerGroup = 2

		Result:
		AI AI Backend Backend DevOps

		NOT:
		AI AI AI AI Backend
	*/
	MaxPerGroup int
}

func (f Fair) Sample(
	query types.Query,
	items []types.ScoredItem,
	k int,
) []types.ScoredItem {

	if len(items) == 0 {
		return nil
	}

	if f.GroupKey == "" {
		return items[:min(k, len(items))]
	}

	if f.MaxPerGroup <= 0 {
		f.MaxPerGroup = 1
	}

	/*
		Sort by score descending.

		Fair sampling works on
		top-ranked candidates first.
	*/
	sorted := make(
		[]types.ScoredItem,
		len(items),
	)

	copy(sorted, items)

	sort.Slice(
		sorted,
		func(i, j int) bool {
			return sorted[i].Score >
				sorted[j].Score
		},
	)

	result := make(
		[]types.ScoredItem,
		0,
		k,
	)

	groupCounts := map[string]int{}

	/*
		First pass:
		apply fairness constraints.
	*/
	for _, item := range sorted {

		if len(result) >= k {
			break
		}

		group := metadataGroup(
			item,
			f.GroupKey,
		)

		/*
			Skip items exceeding
			group quota.
		*/
		if groupCounts[group] >= f.MaxPerGroup {
			continue
		}

		groupCounts[group]++

		result = append(
			result,
			item,
		)
	}

	/*
		Second pass:
		fill remaining slots
		if fairness constraints
		removed too many items.
	*/
	if len(result) < k {

		exists := map[string]struct{}{}

		for _, item := range result {
			exists[item.Item.ID] = struct{}{}
		}

		for _, item := range sorted {

			if len(result) >= k {
				break
			}

			if _, ok := exists[item.Item.ID]; ok {
				continue
			}

			result = append(
				result,
				item,
			)
		}
	}

	return result
}

/*
	metadataGroup extracts
	group values safely.

	If metadata is missing,
	items fall into "unknown".
*/
func metadataGroup(
	item types.ScoredItem,
	key string,
) string {

	if item.Item.Metadata == nil {
		return "unknown"
	}

	value, ok := item.Item.Metadata[key]

	if !ok {
		return "unknown"
	}

	switch v := value.(type) {

	case string:
		return v

	case int:
		return string(rune(v))

	case int64:
		return string(rune(v))

	case float64:
		return string(rune(int(v)))

	default:
		return "unknown"
	}
}

func min(
	a,
	b int,
) int {

	if a < b {
		return a
	}

	return b
}