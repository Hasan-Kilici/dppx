package native

/*
#cgo CXXFLAGS: -std=c++17 -I../cpp/include
#cgo CFLAGS: -I../cpp/include
#cgo LDFLAGS: -L../cpp/build -ldppx_native -lstdc++

#include "bridge.h"
*/
import "C"

import "unsafe"

func (i *Index) Search(query []float32, k int) []SearchResult {
	if i.ptr == nil || len(query) == 0 || k <= 0 {
		return nil
	}

	results := make([]C.CSearchResult, k)
	C.dppx_search(
		i.ptr,
		(*C.float)(unsafe.Pointer(&query[0])),
		C.int(len(query)),
		C.int(k),
		(*C.CSearchResult)(unsafe.Pointer(&results[0])),
	)

	out := make([]SearchResult, 0, k)
	for _, r := range results {
		out = append(out, SearchResult{
			ID:    uint64(r.id),
			Score: float32(r.score),
		})
	}

	return out
}

func (i *Index) SearchWithEF(query []float32, k int, efSearch int) []SearchResult {
	if i.ptr == nil || len(query) == 0 || k <= 0 {
		return nil
	}

	if efSearch <= 0 {
		efSearch = i.config.EFSearch
		if efSearch <= 0 {
			efSearch = 50
		}
	}

	results := make([]C.CSearchResult, k)
	C.dppx_search_options(
		i.ptr,
		(*C.float)(unsafe.Pointer(&query[0])),
		C.int(len(query)),
		C.int(k),
		C.int(efSearch),
		(*C.CSearchResult)(unsafe.Pointer(&results[0])),
	)

	out := make([]SearchResult, 0, k)
	for _, r := range results {
		out = append(out, SearchResult{
			ID:    uint64(r.id),
			Score: float32(r.score),
		})
	}

	return out
}