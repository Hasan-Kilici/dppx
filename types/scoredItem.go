package types

// ScoredItem couples an item with its computed ranking score.
type ScoredItem struct {
	Item  Item
	Score float64
}
