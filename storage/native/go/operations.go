package native

import (
	"container/heap"
)

// MergeMode controls how scores are combined when merging search results.
type MergeMode int

const (
	// MergeModeMax keeps the highest score for an ID.
	MergeModeMax MergeMode = iota
	// MergeModeSum sums scores across result sets.
	MergeModeSum
	// MergeModeAvg averages scores across result sets.
	MergeModeAvg
)

func mergeSearchScores(mode MergeMode, counts map[uint64]int, scoreSums map[uint64]float32, bestScores map[uint64]float32, id uint64, score float32) {
	counts[id]++
	scoreSums[id] += score

	switch mode {
	case MergeModeMax:
		if existing, ok := bestScores[id]; !ok || score > existing {
			bestScores[id] = score
		}
	case MergeModeSum, MergeModeAvg:
		bestScores[id] = scoreSums[id]
	default:
		if existing, ok := bestScores[id]; !ok || score > existing {
			bestScores[id] = score
		}
	}
}

type resultHeap struct {
	data []SearchResult
	k    int
}

func (h *resultHeap) Len() int {
	return len(h.data)
}

func (h *resultHeap) Less(i, j int) bool {
	return h.data[i].Score < h.data[j].Score
}

func (h *resultHeap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *resultHeap) Push(x any) {
	h.data = append(h.data, x.(SearchResult))
}

func (h *resultHeap) Pop() any {
	old := h.data
	n := len(old)
	x := old[n-1]
	h.data = old[:n-1]
	return x
}

func (h *resultHeap) Add(item SearchResult) {
	if len(h.data) < h.k {
		heap.Push(h, item)
		return
	}

	if item.Score > h.data[0].Score {
		heap.Pop(h)
		heap.Push(h, item)
	}
}

func (h *resultHeap) Result() []SearchResult {
	result := make([]SearchResult, 0, len(h.data))
	for h.Len() > 0 {
		result = append(result, heap.Pop(h).(SearchResult))
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

func newResultHeap(k int) *resultHeap {
	return &resultHeap{
		data: make([]SearchResult, 0, k),
		k:    k,
	}
}

// UnionSearchResults returns a deduplicated top-k set from multiple search result sets.
// Duplicate IDs are collapsed using the provided merge mode.
func UnionSearchResults(mode MergeMode, k int, sets ...[]SearchResult) []SearchResult {
	if k <= 0 {
		return nil
	}

	scoreSums := make(map[uint64]float32)
	bestScores := make(map[uint64]float32)
	counts := make(map[uint64]int)

	for _, set := range sets {
		if set == nil {
			continue
		}

		for _, item := range set {
			mergeSearchScores(mode, counts, scoreSums, bestScores, item.ID, item.Score)
		}
	}

	if len(bestScores) == 0 {
		return nil
	}

	heap := newResultHeap(k)
	for id, score := range bestScores {
		if mode == MergeModeAvg {
			score = score / float32(counts[id])
		}
		heap.Add(SearchResult{ID: id, Score: score})
	}

	return heap.Result()
}

// IntersectSearchResults returns the top-k IDs that appear in every result set.
// Scores are merged using the provided merge mode.
func IntersectSearchResults(mode MergeMode, k int, sets ...[]SearchResult) []SearchResult {
	if k <= 0 || len(sets) == 0 {
		return nil
	}

	scoreSums := make(map[uint64]float32)
	bestScores := make(map[uint64]float32)
	counts := make(map[uint64]int)

	for _, set := range sets {
		if set == nil {
			return nil
		}

		seen := make(map[uint64]struct{}, len(set))
		for _, item := range set {
			if _, already := seen[item.ID]; already {
				continue
			}
			seen[item.ID] = struct{}{}
			mergeSearchScores(mode, counts, scoreSums, bestScores, item.ID, item.Score)
		}
	}

	heap := newResultHeap(k)
	for id, score := range bestScores {
		if counts[id] != len(sets) {
			continue
		}
		if mode == MergeModeAvg {
			score = score / float32(counts[id])
		}
		heap.Add(SearchResult{ID: id, Score: score})
	}

	return heap.Result()
}

const (
	defaultIntersectCandidateMultiplier = 10
	minIntersectCandidateWindow         = 100
)

func intersectCandidateWindow(k int) int {
	if k <= 0 {
		return minIntersectCandidateWindow
	}

	window := k * defaultIntersectCandidateMultiplier
	if window < minIntersectCandidateWindow {
		window = minIntersectCandidateWindow
	}
	return window
}

// SearchUnion executes a sequence of native searches and merges the result sets.
// This is a lightweight prototype helper layered on top of native Search().
func (i *Index) SearchUnion(queries [][]float32, k int, mode MergeMode) []SearchResult {
	if i == nil || i.ptr == nil || len(queries) == 0 || k <= 0 {
		return nil
	}

	sets := make([][]SearchResult, 0, len(queries))
	for _, q := range queries {
		if len(q) == 0 {
			continue
		}
		sets = append(sets, i.Search(q, k))
	}

	return UnionSearchResults(mode, k, sets...)
}

// SearchIntersect executes native searches and returns only IDs found in every query.
// This is a lightweight prototype composition helper; the candidate-window expansion
// is intentionally a temporary heuristic and should evolve to adaptive depth control.
func (i *Index) SearchIntersect(queries [][]float32, k int, mode MergeMode) []SearchResult {
	return i.SearchIntersectWithCandidateWindow(queries, k, mode, intersectCandidateWindow(k))
}

// SearchIntersectWithCandidateWindow performs intersection using an explicit candidate window.
// candidateWindow is an explicit expansion knob; a value of 0 uses the current heuristic.
// This API is intended to support future adaptive retrieval depth without hardwiring a single multiplier.
func (i *Index) SearchIntersectWithCandidateWindow(queries [][]float32, k int, mode MergeMode, candidateWindow int) []SearchResult {
	if i == nil || i.ptr == nil || len(queries) == 0 || k <= 0 {
		return nil
	}

	if candidateWindow <= 0 {
		candidateWindow = intersectCandidateWindow(k)
	}

	sets := make([][]SearchResult, 0, len(queries))
	for _, q := range queries {
		if len(q) == 0 {
			return nil
		}
		sets = append(sets, i.Search(q, candidateWindow))
	}

	return IntersectSearchResults(mode, k, sets...)
}
