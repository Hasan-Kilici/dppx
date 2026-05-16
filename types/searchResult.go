package types

// SearchResult contains an item result with score and similarity metadata.
type SearchResult struct {
	Item       Item
	Score      float64
	Similarity float64
}
