package native

import "unsafe"

type SearchResult struct {
	ID    uint64
	Score float32
}

type Index struct {
	ptr       unsafe.Pointer
	dimension int
	config    Config
}