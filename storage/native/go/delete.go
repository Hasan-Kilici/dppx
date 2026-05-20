package native

/*
#cgo CXXFLAGS: -std=c++17 -I../cpp/include
#cgo CFLAGS: -I../cpp/include
#cgo LDFLAGS: -L../cpp/build -ldppx_native -lstdc++

#include "bridge.h"
*/
import "C"

func (i *Index) Delete(
	id uint64,
) {

	C.dppx_delete(
		i.ptr,
		C.ulonglong(id),
	)
}