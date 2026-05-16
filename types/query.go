package types

// Query is the input context for search operations, including vector and metadata.
type Query struct {
	UserID   string
	Vector   Vector
	Norm     float32
	Metadata map[string]any
}
