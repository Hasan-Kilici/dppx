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

func New(cfg Config) *Index {
	path := cfg.Path
	if path == "" {
		path = "."
	}

	cfgPtr := C.CString(path)
	defer C.free(unsafe.Pointer(cfgPtr))

	config := C.NativeConfig{
		dimension:      C.int(cfg.Dimension),
		path:           cfgPtr,
		hnswM:          C.int(cfg.HNSWM),
		efConstruction: C.int(cfg.EFConstruction),
		efSearch:       C.int(cfg.EFSearch),
		maxSegmentSize: C.int(cfg.MaxSegmentSize),
		useHNSW:        C.int(boolToInt(cfg.UseHNSW)),
		enableWAL:      C.int(boolToInt(cfg.EnableWAL)),
		storePayload:   C.int(boolToInt(cfg.StorePayload)),
	}

	ptr := C.dppx_create_index(config)

	return &Index{
		ptr:       ptr,
		dimension: cfg.Dimension,
		config:    cfg,
	}
}

func (i *Index) Close() {
	if i.ptr != nil {
		C.dppx_free_index(i.ptr)
		i.ptr = nil
	}
}
func (i *Index) Flush() {
	if i.ptr == nil {
		return
	}

	C.dppx_flush(i.ptr)
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}