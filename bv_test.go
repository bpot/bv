package bv

import (
	"math/rand"
	"testing"
	"time"
)

func TestBVMultipleOf64(t *testing.T) {
	bv := New(64 * 64)
	var i uint64
	for i = 0; i < 64*64; i++ {
		if bv.Get(i) {
			t.Errorf("bit %d should be unset", i)
		}
		bv.Set(i, true)
	}
}

func TestBVGetSet(t *testing.T) {
	bv := New(1024)
	var i uint64
	for i = 0; i < 1024; i++ {
		if bv.Get(i) {
			t.Errorf("bit %d should be unset", i)
		}
		bv.Set(i, true)
	}

	for i = 0; i < 1024; i++ {
		if !bv.Get(i) {
			t.Errorf("bit %d should be set", i)
		}
		bv.Set(i, false)
	}

	for i = 0; i < 1024; i++ {
		if bv.Get(i) {
			t.Errorf("bit %d should be unset", i)
		}
	}
}

const randMax = 1 << 20

func TestRandSetGet(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	bv := New(randMax)
	for i := 0; i < 10000; i++ {
		v := uint64(rnd.Intn(randMax))
		bv.Set(v, true)
		if !bv.Get(v) {
			t.Fatalf("expected %d to be set", v)
		}
		bv.Set(v, false)
		if bv.Get(v) {
			t.Fatalf("expected %d to not be set", v)
		}
	}
}

func TestSetGetInt(t *testing.T) {
	examples := []struct {
		offset uint
		bitLen uint8
		v      uint64
	}{
		{633431, 12, 3461},
		{629854, 49, 298629597710341},
	}

	bv := New(randMax)
	for i := 0; i < 999999; i++ {
		bv.Set(uint64(i), true)
	}
	for _, ex := range examples {
		bv.SetInt(int(ex.offset), ex.bitLen, ex.v)
		actualV := bv.GetInt(ex.offset, ex.bitLen)
		if ex.v != actualV {
			t.Errorf("expected %d; got %d", ex.v, actualV)
		}
	}
}

func TestRandSetGetInt(t *testing.T) {
	bv := New(randMax)

	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 1000000; i++ {
		// Randomly generate an offset, an int length and an int which fits in int length.
		bitLen := uint8(rnd.Intn(63) + 1)
		offset := rnd.Intn(randMax - int(bitLen))
		vExpected := uint64(rnd.Intn(1<<uint(bitLen) - 1))

		bv.SetInt(offset, bitLen, vExpected)
		vActual := bv.GetInt(uint(offset), bitLen)
		if vExpected != vActual {
			t.Fatalf("SetInt/GetInt failed: offset %d, bitLen %d, expectedV %d, gotV %d", offset, bitLen, vExpected, vActual)
		}
	}
}

type ints struct {
	offset int
	width  uint8
	v      uint64
}

func TestRandInt(t *testing.T) {
	bv := New(randMax)

	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	var offset int
	var is []ints
	for i := 0; i < 10000; i++ {
		bitLen := uint8(rnd.Intn(63) + 1)
		vExpected := uint64(rnd.Intn(1<<uint(bitLen) - 1))
		bv.SetInt(offset, bitLen, vExpected)

		is = append(is, ints{
			offset: offset,
			width:  bitLen,
			v:      vExpected,
		})

		offset += int(bitLen)
	}

	for _, i := range is {
		actualV := bv.GetInt(uint(i.offset), i.width)
		if i.v != actualV {
			t.Errorf("expected %d;got %d", i.v, actualV)
		}
	}
}
