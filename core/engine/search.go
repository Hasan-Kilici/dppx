package engine

import (
	"context"
	"runtime"
	"sync"

	"github.com/hasan-kilici/dppx/core/topk"
	"github.com/hasan-kilici/dppx/types"
)

func (e *Engine) Search(
	ctx context.Context,
	query types.Query,
	k int,
) ([]types.ScoredItem, error) {

	/*Resolve candidate pool size.*/
	candidatePool := e.config.CandidatePool

	if candidatePool <= 0 {
		candidatePool = k
	}

	if candidatePool < k {
		candidatePool = k
	}

	/*Retrieve ANN candidates.*/

	items, err := e.config.Retriever.Search(
		ctx,
		query,
		candidatePool,
	)

	if err != nil {
		return nil, err
	}

	jobs := make(
		chan types.Item,
		len(items),
	)

	workerCount := runtime.NumCPU()

	localHeaps := make(
		[]*topk.MinHeap,
		workerCount,
	)

	var wg sync.WaitGroup

	/* Parallel scoring workers.*/
	for i := 0; i < workerCount; i++ {
		localHeaps[i] = topk.New(k)
		wg.Add(1)
		go func(
			localHeap *topk.MinHeap,
		) {
			defer wg.Done()
			for item := range jobs {
				similarityScore := e.config.Similarity(
					query.Vector,
					item.Vector,
				)

				customScore := 0.0
				if e.config.Scoring != nil {
					customScore =
						e.config.Scoring(
							query,
							item,
						)
				}
				/*
					Temporary score fusion.

					TODO:
					make configurable.
				*/

				finalScore :=
					(similarityScore * 0.7) +
						(customScore * 0.3)

				localHeap.Add(
					types.ScoredItem{
						Item:  item,
						Score: finalScore,
					},
				)
			}

		}(localHeaps[i])
	}

	/*Dispatch jobs.*/
	for _, item := range items {
		jobs <- item
	}

	close(jobs)

	wg.Wait()

	/*Merge local heaps.*/
	globalHeap := topk.New(k)
	for _, h := range localHeaps {
		for _, item := range h.Result() {
			globalHeap.Add(
				item,
			)
		}
	}

	scored := globalHeap.Result()
	/*Optional diversification stage.*/
	if e.config.Sampler != nil {
		scored =
			e.config.Sampler.Sample(
				query,
				scored,
				k,
			)
	}

	return scored, nil
}