package topk

import (
	"container/heap"
	"github.com/hasan-kilici/dppx/types"
)

// MinHeap maintains the top-k scored items using min-heap semantics.
// The root always holds the smallest score, so lower-scoring items are discarded first.
type MinHeap struct {
	data []types.ScoredItem
	k    int
}

func New(k int) *MinHeap {
	return &MinHeap{
		data: make([]types.ScoredItem, 0, k),
		k:    k,
	}
}

func (h *MinHeap) Len() int {
	return len(h.data)
}

func (h *MinHeap) Less(i, j int) bool {
	return h.data[i].Score < h.data[j].Score
}

func (h *MinHeap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *MinHeap) Push(x any) {
	h.data = append(h.data, x.(types.ScoredItem))
}

func (h *MinHeap) Pop() any {
	old := h.data
	n := len(old)
	x := old[n-1]
	h.data = old[:n-1]
	return x
}

func (h *MinHeap) Add(item types.ScoredItem) {
	if len(h.data) < h.k {
		heap.Push(h, item)
		return
	}

	// Only replace the root if the incoming score is better than the current minimum.
	if item.Score > h.data[0].Score {
		heap.Pop(h)
		heap.Push(h, item)
	}
}

// Result empties the heap and returns results sorted from highest to lowest score.
func (h *MinHeap) Result() []types.ScoredItem {
	result := make([]types.ScoredItem, 0, len(h.data))

	for h.Len() > 0 {
		result = append(result, heap.Pop(h).(types.ScoredItem))
	}

	// reverse (max → min)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}
