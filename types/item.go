package types

// Item represents a candidate object with vector data and optional metadata.
type Item struct {
	ID       string
	Vector   Vector
	Norm     float32
	Metadata map[string]any
}
