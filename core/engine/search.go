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

	items, err := e.config.Retriever.Search(
		ctx,
		query,
		e.config.CandidatePool,
	)

	if err != nil {
		return nil, err
	}

	jobs := make(chan types.Item, len(items))

	workerCount := runtime.NumCPU()

	localHeaps := make([]*topk.MinHeap, workerCount)

	var wg sync.WaitGroup

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
					query.Norm,
					item.Norm,
				)

				customScore := 0.0

				if e.config.Scoring != nil {

					customScore = e.config.Scoring(
						query,
						item,
					)
				}

				finalScore :=
					(similarityScore * 0.7) +
						(customScore * 0.3)

				localHeap.Add(types.ScoredItem{
					Item:  item,
					Score: finalScore,
				})
			}

		}(localHeaps[i])
	}

	for _, item := range items {
		jobs <- item
	}

	close(jobs)

	wg.Wait()

	globalHeap := topk.New(k)

	for _, h := range localHeaps {

		for _, item := range h.Result() {
			globalHeap.Add(item)
		}
	}

	scored := globalHeap.Result()

	if e.config.Sampler != nil {

		scored = e.config.Sampler.Sample(
			query,
			scored,
			k,
		)
	}

	return scored, nil
}