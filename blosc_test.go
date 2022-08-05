package blosc_test

import (
	"testing"
	"unsafe"

	blosc "github.com/seerai/go-blosc"
)

func TestRoundTrip(t *testing.T) {

	buf := make([]int64, 10000)
	for i := 0; i < len(buf); i++ {
		buf[i] = int64(112)
	}

	err := blosc.SetCompressor("lz4")
	if err != nil {
		t.Error(err)
	}

	cmp, err := blosc.Compress(1, true, buf)
	if err != nil {
		t.Error(err)
	}
	dec := blosc.Decompress(cmp)

	if len(dec)/int(unsafe.Sizeof(buf[0])) != len(buf) {
		t.Fatal("unexpected length on decompression")
	}

	obuf := (*(*[1 << 32]int64)(unsafe.Pointer(&dec[0])))[:len(dec)/int(unsafe.Sizeof(buf[0]))]
	for i, o := range obuf {
		if o != buf[i] {
			t.Fatal("unequal after decompression")
		}
	}

}
