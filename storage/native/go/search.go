package native
/*
#cgo CXXFLAGS: -std=c++17 -I../cpp/include
#cgo CFLAGS: -I../cpp/include
#cgo LDFLAGS: -L../cpp/build -ldppx_native -lstdc++

#include <stdlib.h>
#include "bridge.h"
*/
import "C"

import "unsafe"

type SearchCursor struct {
	ptr C.CSearchCursor
}

func (c *SearchCursor) Next() (SearchResult, bool) {
	if c == nil || c.ptr == nil {
		return SearchResult{}, false
	}

	var result C.CSearchResult
	ok := int(C.dppx_search_cursor_next(c.ptr, &result))
	if ok == 0 {
		return SearchResult{}, false
	}

	return SearchResult{
		ID:    uint64(result.id),
		Score: float32(result.score),
	}, true
}

func (c *SearchCursor) Close() {
	if c == nil || c.ptr == nil {
		return
	}
	C.dppx_search_cursor_free(c.ptr)
	c.ptr = nil
}

func (i *Index) Search(query []float32, k int) []SearchResult {
	return i.searchWithCount(query, k, 0, false)
}

func (i *Index) SearchCursor(query []float32, k int) *SearchCursor {
	if i == nil || i.ptr == nil || len(query) == 0 || k <= 0 {
		return nil
	}

	cursor := C.dppx_search_cursor(
		i.ptr,
		(*C.float)(unsafe.Pointer(&query[0])),
		C.int(len(query)),
		C.int(k),
	)
	if cursor == nil {
		return nil
	}

	return &SearchCursor{ptr: cursor}
}

type SearchOptions struct {
	EFSearch int
}

func (i *Index) SearchWithOptions(query []float32, k int, opts SearchOptions) []SearchResult {
	if i.ptr == nil || len(query) == 0 || k <= 0 {
		return nil
	}

	efSearch := opts.EFSearch
	if efSearch <= 0 {
		efSearch = i.config.EFSearch
		if efSearch <= 0 {
			efSearch = 50
		}
	}

	return i.SearchWithEF(query, k, efSearch)
}

func (i *Index) SearchCursorWithEF(query []float32, k int, efSearch int) *SearchCursor {
	if i == nil || i.ptr == nil || len(query) == 0 || k <= 0 {
		return nil
	}

	if efSearch <= 0 {
		efSearch = i.config.EFSearch
		if efSearch <= 0 {
			efSearch = 50
		}
	}

	cursor := C.dppx_search_cursor_options(
		i.ptr,
		(*C.float)(unsafe.Pointer(&query[0])),
		C.int(len(query)),
		C.int(k),
		C.int(efSearch),
	)
	if cursor == nil {
		return nil
	}

	return &SearchCursor{ptr: cursor}
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

	return i.searchWithCount(query, k, efSearch, true)
}

func (i *Index) searchWithCount(query []float32, k int, efSearch int, useOptions bool) []SearchResult {
	if i.ptr == nil || len(query) == 0 || k <= 0 {
		return nil
	}

	storageSize := C.size_t(k) * C.size_t(unsafe.Sizeof(C.CSearchResult{}))
	raw := C.malloc(storageSize)
	if raw == nil {
		return nil
	}
	defer C.free(raw)

	count := 0
	if useOptions {
		count = int(C.dppx_search_options(
			i.ptr,
			(*C.float)(unsafe.Pointer(&query[0])),
			C.int(len(query)),
			C.int(k),
			C.int(efSearch),
			(*C.CSearchResult)(raw),
		))
	} else {
		count = int(C.dppx_search(
			i.ptr,
			(*C.float)(unsafe.Pointer(&query[0])),
			C.int(len(query)),
			C.int(k),
			(*C.CSearchResult)(raw),
		))
	}

	if count <= 0 {
		return nil
	}

	results := unsafe.Slice((*C.CSearchResult)(raw), count)
	out := make([]SearchResult, count)
	for i, r := range results {
		out[i] = SearchResult{
			ID:    uint64(r.id),
			Score: float32(r.score),
		}
	}

	return out
}
