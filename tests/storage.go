package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	native "github.com/hasan-kilici/dppx/storage/native/go"
)

func main() {

	tmpDir := "./dppx_data"

	_ = os.MkdirAll(
		tmpDir,
		0755,
	)

	cfg := native.Config{
		Dimension: 128,

		Path: tmpDir,

		UseHNSW: true,

		HNSWM: 16,

		EFConstruction: 200,

		EFSearch: 64,

		MaxSegmentSize: 4 * 1024 * 1024,

		StorePayload: false,

		EnableWAL: true,
	}

	fmt.Println("===================================")
	fmt.Println("DPPX Native Storage Demo")
	fmt.Println("===================================")

	startCreate := time.Now()

	idx := native.New(cfg)

	if idx == nil {
		panic("failed to create index")
	}

	defer idx.Close()

	fmt.Printf(
		"index created in %s\n",
		time.Since(startCreate),
	)

	/*
		Insert synthetic vectors
	*/

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

	/*
		Search benchmark
	*/

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

	searchLatency := time.Since(
		searchStart,
	)

	fmt.Printf(
		"search latency: %s\n",
		searchLatency,
	)

	fmt.Println()
	fmt.Println("========== TOP RESULTS ==========")

	for i, r := range results {

		fmt.Printf(
			"%d. ID=%d SCORE=%.6f\n",
			i+1,
			r.ID,
			r.Score,
		)
	}

	/*
		SearchWithEF benchmark
	*/

	fmt.Println()
	fmt.Println("running efSearch benchmark...")

	efStart := time.Now()

	resultsEF := idx.SearchWithEF(
		query,
		10,
		128,
	)

	efLatency := time.Since(
		efStart,
	)

	fmt.Printf(
		"efSearch latency: %s\n",
		efLatency,
	)

	fmt.Println()
	fmt.Println("========== EF SEARCH RESULTS ==========")

	for i, r := range resultsEF {

		fmt.Printf(
			"%d. ID=%d SCORE=%.6f\n",
			i+1,
			r.ID,
			r.Score,
		)
	}

	/*
		Delete test
	*/

	fmt.Println()
	fmt.Println("testing delete...")

	idx.Delete(1)

	fmt.Println("deleted vector id=1")

	/*
		Flush test
	*/

	fmt.Println()
	fmt.Println("flushing storage...")

	flushStart := time.Now()

	idx.Flush()

	fmt.Printf(
		"flush completed in %s\n",
		time.Since(flushStart),
	)

	/*
		Check WAL
	*/

	walPath := filepath.Join(
		tmpDir,
		"wal",
		"wal.log",
	)

	if info, err := os.Stat(
		walPath,
	); err == nil {

		fmt.Printf(
			"WAL size: %d bytes\n",
			info.Size(),
		)
	}

	/*
		Check segments
	*/

	segmentPath := filepath.Join(
		tmpDir,
		"segments",
	)

	entries, err := os.ReadDir(
		segmentPath,
	)

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

func makeVector(
	dim int,
	seed float32,
) []float32 {

	vector := make(
		[]float32,
		dim,
	)

	for i := range vector {

		vector[i] =
			float32(i)*0.001 +
				seed*0.0001
	}

	return vector
}