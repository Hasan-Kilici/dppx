package engine

import (
	"runtime"
	"sync"

	"github.com/hasan-kilici/dppx/core/topk"
	"github.com/hasan-kilici/dppx/types"
)

// Search ranks items for a query in parallel and returns the top-k scored results.
// It uses a worker-per-CPU model with local heaps to minimize heap contention.
func (e *Engine) Search(
	query types.Query,
	items []types.Item,
	k int,
) []types.ScoredItem {

	workerCount := runtime.NumCPU()

	jobs := make(chan types.Item, len(items))

	// Each worker maintains its own min-heap to avoid synchronization overhead during scoring.
	localHeaps := make([]*topk.MinHeap, workerCount)

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		localHeaps[i] = topk.New(k)
	}

	// Launch one scorer goroutine per CPU core. Each worker consumes items from the job channel.
	for i := 0; i < workerCount; i++ {

		wg.Add(1)
		heap := localHeaps[i]

		go func() {
			defer wg.Done()

			for item := range jobs {

				sim := e.config.Similarity(
					query.Vector,
					item.Vector,
					query.Norm,
					item.Norm,
				)

				custom := 0.0
				if e.config.Scoring != nil {
					custom = e.config.Scoring(query, item)
				}

				// Blend similarity and custom business score into a single ranking value.
				final := (sim * 0.7) + (custom * 0.3)

				heap.Add(types.ScoredItem{
					Item:  item,
					Score: final,
				})
			}
		}()
	}

	// Feed all items into the shared job channel, then close it to signal completion.
	go func() {
		for _, item := range items {
			jobs <- item
		}
		close(jobs)
	}()

	wg.Wait()

	// Merge phase combines local worker heaps into a final global top-k result.
	// This preserves parallel throughput while producing a single ranked output.
	finalHeap := topk.New(k)

	for _, h := range localHeaps {

		for _, item := range h.Result() {
			finalHeap.Add(item)
		}
	}

	return finalHeap.Result()
}
