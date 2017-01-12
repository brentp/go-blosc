package blosc

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestRoundTrip(t *testing.T) {

	buf := make([]int64, 10000)
	for i := 0; i < len(buf); i++ {
		buf[i] = int64(112)
	}

	cmp := Compress(1, true, buf)
	dec := Decompress(cmp)

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

func TestUint16(t *testing.T) {

	tt := []uint16{1, 2, 3, 4, 5, 6, 7, 8}

	cmp := Compress(1, true, tt)
	dec := Decompress(cmp).Uint16s()
	if len(dec)/int(unsafe.Sizeof(tt[0])) != len(tt) {
		t.Fatal("unexpected length on decompression")
	}
	if !reflect.DeepEqual(dec, tt) {
		t.Fatalf("unequal after decompression %+v, %+v", dec, tt)
	}

}
