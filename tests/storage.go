package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	native "github.com/hasan-kilici/dppx/storage/native/go"
)

func main() {

	// Base directory used for WAL and segment persistence.
	tmpDir := "./dppx_data"

	_ = os.MkdirAll(tmpDir, 0755)

	cfg := native.Config{
		Dimension: 128,
		Path:      tmpDir,

		// Enables approximate nearest neighbor indexing using HNSW.
		UseHNSW: true,

		// HNSW graph connectivity parameter.
		// Higher values generally improve recall at the cost of memory.
		HNSWM: 16,

		// Controls graph build quality during indexing.
		// Higher values improve accuracy but slow insertion time.
		EFConstruction: 200,

		// Runtime search exploration factor.
		// Higher values improve recall but increase latency.
		EFSearch: 64,

		// Maximum on-disk segment size before rolling into a new segment.
		MaxSegmentSize: 4 * 1024 * 1024, // 4 MB

		// Payload storage is disabled for this benchmark-focused example.
		StorePayload: false,

		// Enables write-ahead logging for durability/recovery.
		EnableWAL: true,
	}

	fmt.Println("===================================")
	fmt.Println("DPPX Native Storage Demo")
	fmt.Println("===================================")

	idx := native.New(cfg)
	if idx == nil {
		panic("failed to create index")
	}
	defer idx.Close()

	fmt.Println("index created")

	fmt.Println()
	fmt.Println("inserting vectors...")

	insertStart := time.Now()

	total := 10000

	for i := 0; i < total; i++ {

		vector := makeVector(
			cfg.Dimension,
			float32(i),
		)

		idx.Insert(
			uint64(i+1),
			vector,
		)
	}

	fmt.Printf(
		"inserted %d vectors in %s\n",
		total,
		time.Since(insertStart),
	)

	query := makeVector(
		cfg.Dimension,
		42,
	)

	fmt.Println()
	fmt.Println("running search benchmark...")

	searchStart := time.Now()

	results := idx.Search(
		query,
		10,
	)

	fmt.Printf(
		"search latency: %s\n",
		time.Since(searchStart),
	)

	fmt.Println()
	fmt.Printf("requested=10 actual=%d\n", len(results))
	fmt.Println("========== TOP RESULTS ==========")
	printSearchResults(results)
	fmt.Println()
	fmt.Println("running efSearch benchmark...")

	efStart := time.Now()

	// Uses a larger EF value than the default runtime configuration
	// to measure recall/latency tradeoffs dynamically.
	resultsEF := idx.SearchWithEF(
		query,
		10,
		128,
	)

	fmt.Printf(
		"efSearch latency: %s\n",
		time.Since(efStart),
	)

	fmt.Println()
	fmt.Printf("requested=10 actual=%d\n", len(resultsEF))
	fmt.Println("========== EF SEARCH RESULTS ==========")
	printSearchResults(resultsEF)
	fmt.Println()
	fmt.Println("running explicitly configured search options...")

	options := native.SearchOptions{
		EFSearch: 128,
	}

	resultsOpts := idx.SearchWithOptions(
		query,
		10,
		options,
	)

	fmt.Println("========== EXPLICIT SEARCH OPTIONS RESULTS ==========")
	fmt.Printf("requested=10 actual=%d\n", len(resultsOpts))
	printSearchResults(resultsOpts)
	fmt.Println()
	fmt.Println("running cursor-based search...")

	cursor := idx.SearchCursor(query, 10)
	if cursor != nil {
		for i := 0; ; i++ {
			result, ok := cursor.Next()
			if !ok {
				break
			}
			fmt.Printf("%d. ID=%d SCORE=%.6f\n", i+1, result.ID, result.Score)
		}
		cursor.Close()
	}

	fmt.Println()
	fmt.Println("running cursor-based search with explicit EF...")

	cursorEF := idx.SearchCursorWithEF(query, 10, 128)
	if cursorEF != nil {
		for i := 0; ; i++ {
			result, ok := cursorEF.Next()
			if !ok {
				break
			}
			fmt.Printf("%d. ID=%d SCORE=%.6f\n", i+1, result.ID, result.Score)
		}
		cursorEF.Close()
	}

	fmt.Println()
	fmt.Println("testing delete semantics...")

	idx.Delete(1)

	fmt.Println("deleted vector id=1")

	resultsAfterDelete := idx.Search(
		query,
		10,
	)

	fmt.Println("========== RESULTS AFTER DELETE ==========")

	fmt.Printf("requested=10 actual=%d\n", len(resultsAfterDelete))
	printSearchResults(resultsAfterDelete)

	fmt.Println()
	fmt.Println("testing union and intersect operations...")

	query2 := makeVector(
		cfg.Dimension,
		84,
	)

	// Union merges candidate sets from multiple searches.
	// MergeModeMax keeps the highest score per document.
	union := idx.SearchUnion(
		[][]float32{
			query,
			query2,
		},
		10,
		native.MergeModeMax,
	)

	fmt.Println("========== UNION RESULTS ==========")
	fmt.Printf("requested=10 actual=%d\n", len(union))
	printSearchResults(union)

	// Intersect keeps only shared candidates between searches.
	// MergeModeAvg averages scores across matching hits.
	intersect := idx.SearchIntersect(
		[][]float32{
			query,
			query2,
		},
		10,
		native.MergeModeAvg,
	)

	fmt.Println("========== INTERSECT RESULTS ==========")
	fmt.Printf("requested=10 actual=%d\n", len(intersect))
	printSearchResults(intersect)

	explicitIntersect := idx.SearchIntersectWithCandidateWindow(
		[][]float32{
			query,
			query2,
		},
		10,
		native.MergeModeAvg,
		100,
	)

	fmt.Println("========== INTERSECT WITH CANDIDATE WINDOW RESULTS ==========")
	fmt.Printf("requested=10 actual=%d\n", len(explicitIntersect))
	printSearchResults(explicitIntersect)

	fmt.Println()
	fmt.Println("flushing storage...")

	flushStart := time.Now()

	// Forces pending WAL/segment data to disk.
	idx.Flush()

	fmt.Printf(
		"flush completed in %s\n",
		time.Since(flushStart),
	)

	walPath := filepath.Join(
		tmpDir,
		"wal",
		"wal.log",
	)

	if info, err := os.Stat(walPath); err == nil {

		fmt.Printf(
			"WAL size: %d bytes\n",
			info.Size(),
		)
	}

	segmentPath := filepath.Join(
		tmpDir,
		"segments",
	)

	entries, err := os.ReadDir(segmentPath)

	if err == nil {

		fmt.Println()
		fmt.Println("========== SEGMENTS ==========")

		for _, e := range entries {

			fmt.Printf(
				"- %s\n",
				e.Name(),
			)
		}
	}

	fmt.Println()
	fmt.Println("===================================")
	fmt.Println("DPPX Native Storage Completed")
	fmt.Println("===================================")
}

func printSearchResults(results []native.SearchResult) {
	for i, r := range results {
		fmt.Printf(
			"%d. ID=%d SCORE=%.6f\n",
			i+1,
			r.ID,
			r.Score,
		)
	}
}

func makeVector(
	dim int,
	seed float32,
) []float32 {

	vector := make(
		[]float32,
		dim,
	)

	for i := range vector {

		// Generates deterministic synthetic vectors for benchmarking.
		// Small incremental values help create stable similarity behavior.
		vector[i] =
			float32(i)*0.001 +
				seed*0.0001
	}

	return vector
}