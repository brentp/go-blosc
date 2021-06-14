// Package blosc wraps blosc for compressing numbers.
package blosc

/*
#cgo CFLAGS: -O2 -msse2 -I${SRCDIR}/c-blosc/blosc/ -I${SRCDIR}/c-blosc/internal-complibs/lz4-1.9.2 -D HAVE_LZ4=1
#cgo LDFLAGS: -lpthread
#include "blosc_include.h"

int blosc_set_compressor_wrapper(char* compname) {
	return blosc_set_compressor(compname);
}
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/pkg/errors"
)

func init() {
	C.blosc_init()
}

// SetCompressor sets the compressor that blosc will use
func SetCompressor(name string) error {
	switch name {
	case "blosclz", "lz4", "lz4hc":
		break
	default:
		return errors.New("only lz4, blosclz, lz4hc are currently supported")
	}
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))

	res := C.blosc_set_compressor_wrapper(str)
	if res < 0 {
		return errors.New(fmt.Sprintf("invalid compressor '%s'", name))
	}
	return nil
}

// SetBlocksize for the compressor
func SetBlocksize(blocksize int) {
	C.blosc_set_blocksize(C.ulong(blocksize))
}

// Compress takes a slice of numbers and compresses according to level and shuffle.
func Compress(level int, shuffle bool, slice interface{}) ([]byte, error) {

	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return nil, errors.New("input data must be a slice")
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

	return compressed[:csize], nil
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
