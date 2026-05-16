package similarity

// Func defines a similarity function between two dense vectors.
type Func func(query, item []float64) float64