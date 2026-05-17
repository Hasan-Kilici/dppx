package fusion

type Weighted struct {
	SimilarityWeight float32
	CustomWeight     float32
}

func (w Weighted) Combine(
	s Scores,
) float32 {

	return (s.Similarity * w.SimilarityWeight) +
		(s.Custom * w.CustomWeight)
}