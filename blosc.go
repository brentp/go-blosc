// Package blosc wraps blosc for compressing numbers.
package blosc

/*
#cgo CFLAGS: -O2 -msse2 -I${SRCDIR}/c-blosc/blosc/
#cgo LDFLAGS: -lpthread
#include "blosc_include.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

func init() {
	C.blosc_init()
}

// Compress takes a slice of numbers and compresses according to level and shuffle.
func Compress(level int, shuffle bool, slice interface{}) []byte {

	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		panic("blosc: expected slice to Compress")
	}
	l := rv.Len()
	size := int(rv.Index(0).Type().Size())
	ptr := unsafe.Pointer(rv.Pointer())

	s := 1
	if !shuffle {
		s = 0
	}
	compressed := make([]byte, l*size+C.BLOSC_MAX_OVERHEAD)
	csize := C.blosc_compress(C.int(level), C.int(s), C.size_t(size),
		C.size_t(l*size),
		ptr,
		unsafe.Pointer(&compressed[0]),
		C.size_t(len(compressed)))
	return compressed[:csize]
}

// Decompress takes a byte of compressed data and returns the uncompressed data.
func Decompress(compressed []byte) []byte {

	nbytes := C.size_t(0)
	cbytes := C.size_t(0)
	blksz := C.size_t(0)

	C.blosc_cbuffer_sizes(unsafe.Pointer(&compressed[0]), &nbytes, &cbytes, &blksz)

	data := make([]byte, int(nbytes))
	C.blosc_decompress(unsafe.Pointer(&compressed[0]), unsafe.Pointer(&data[0]), nbytes)
	return data
}
