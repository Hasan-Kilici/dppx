package mongodb

type Config struct {

	// Mongo connection URI.
	//
	// Examples:
	// mongodb://localhost:27017
	// mongodb+srv://...
	URI string

	// Target database name.
	Database string

	// Target collection name.
	Collection string

	// Embedding field name.
	//
	// Example:
	// "embedding"
	VectorField string

	// Atlas Vector Search index name.
	//
	// Only used when:
	// UseAtlasVectorSearch = true
	IndexName string

	// Enables Atlas native vector search.
	//
	// false:
	// local Mongo retrieval mode
	//
	// true:
	// Atlas Vector Search mode
	UseAtlasVectorSearch bool
}