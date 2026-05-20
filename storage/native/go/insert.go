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

func (i *Index) Insert(id uint64, vector []float32) {
	if i.ptr == nil || len(vector) == 0 {
		return
	}

	C.dppx_insert(
		i.ptr,
		C.ulonglong(id),
		(*C.float)(unsafe.Pointer(&vector[0])),
		C.int(len(vector)),
	)
}

func (i *Index) InsertBatch(ids []uint64, vectors []float32, dim int) {
	if i.ptr == nil || len(ids) == 0 || len(vectors) == 0 || dim <= 0 {
		return
	}

	expected := len(ids) * dim
	if len(vectors) != expected {
		return
	}

	C.dppx_insert_batch(
		i.ptr,
		(*C.ulonglong)(unsafe.Pointer(&ids[0])),
		(*C.float)(unsafe.Pointer(&vectors[0])),
		C.int(len(ids)),
		C.int(dim),
	)
}