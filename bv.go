package bv

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

// BV is a bitvector
type BV struct {
	size uint64
	bits []uint64
}

// New returns a bitvector containing size bits
func New(size int) *BV {
	blocks := (size + 63) / 64
	return &BV{
		size: uint64(size),
		bits: make([]uint64, blocks),
	}
}

// NewByteBacked returns a bitvector backed by buf
func NewByteBacked(buf []byte) (bitVector *BV, serializedLength uint64, err error) {
	bv := &BV{}

	bv.size = binary.LittleEndian.Uint64(buf)
	buf = buf[8:]

	blocks := (bv.size + 63) / 64
	if blocks*8 > uint64(len(buf)) {
		return nil, 0, fmt.Errorf("bv: serialized buffer is malformed. too short.")
	}

	bv.bits = byteSliceAsUint64Slice(buf[:blocks*8])
	return bv, 8 + blocks*8, nil
}

// Size returns the number of bits the vector represents.
func (bv *BV) Size() int {
	return int(bv.size)
}

// Reset clears the bit vector.
func (bv *BV) Reset() {
	for i := range bv.bits {
		bv.bits[i] = 0
	}
}

// Get returns true if the ith bit is set.
func (bv *BV) Get(i uint64) bool {
	return (bv.bits[i/64]>>(i%64))&1 == 1
}

// Set sets or clears the ith bit.
func (bv *BV) Set(i uint64, b bool) {
	if b {
		bv.bits[i/64] |= (1 << (i % 64))
	} else {
		bv.bits[i/64] &^= (1 << (i % 64))
	}
}

// GetInt returns an integer representation of bitLen bits at offset
func (bv *BV) GetInt(offset uint, bitLen uint8) uint64 {
	block1 := offset / 64
	block1Offset := offset % 64

	v := bv.bits[block1] >> block1Offset
	if 64-block1Offset < uint(bitLen) {
		v |= (bv.bits[block1+1] << (64 - block1Offset))
	}
	return v & masks[bitLen]
}

// SetInt sets bitLen bits at offset
func (bv *BV) SetInt(offset int, bitLen uint8, v uint64) {
	if offset+int(bitLen) > int(bv.size) {
		panic(fmt.Sprintf("offset+bitLen too large: %d > %d", offset+int(bitLen), bv.size))
	}
	block1 := offset / 64
	block1Offset := uint8(offset % 64)

	block1Mask := uint64(((1 << bitLen) - 1) << block1Offset)
	block1V := (v << block1Offset)
	bv.bits[block1] = (bv.bits[block1] & ^block1Mask) | block1V

	bitsForBlock2 := int(bitLen) - (64 - int(block1Offset))
	if bitsForBlock2 > 0 {
		block2Mask := uint64((1 << bitLen) - 1)
		block2V := v >> (64 - block1Offset)
		bv.bits[block1+1] = (bv.bits[block1+1] & ^block2Mask) | block2V
	}
}

// Equals returns true if the two bit vectors are equal.
func (bv *BV) Equals(bv2 *BV) bool {
	if len(bv.bits) != len(bv2.bits) {
		return false
	}
	for i, w := range bv.bits {
		if w != bv2.bits[i] {
			return false
		}
	}
	return true
}

// SizeInBytes returns the size of the bit vector in bytes.
func (bv *BV) SizeInBytes() uint64 {
	return uint64(len(bv.bits)*8 + 8)
}

// String returns a string representation of the bit vector.
func (bv *BV) String() string {
	b := []byte{'{'}
	for i := 0; i < int(bv.size); i++ {
		if bv.Get(uint64(i)) {
			b = append(b, '1')
		} else {
			b = append(b, '0')
		}
		if i != int(bv.size)-1 {
			b = append(b, []byte(", ")...)
		}
	}
	b = append(b, '}')
	return string(b)
}

// WriteTo serializes the bv to w. It returns the number of bits written.
func (bv *BV) WriteTo(w io.Writer) (n int, err error) {
	// Write size of the bitvector
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(bv.size))
	_, err = w.Write(buf)
	if err != nil {
		return 0, err
	}

	return w.Write(uint64SliceAsByteSlice(bv.bits))
}

// SerializedSize returns the size of the bv when serialized. i.e. WriteTo
func (bv *BV) SerializedSize() int {
	return 8 * len(bv.bits) / 8
}

func (bv *BV) readFrom(r io.Reader) (n int, err error) {
	err = binary.Read(r, binary.LittleEndian, &bv.size)
	if err != nil {
		return 0, err
	}
	bv.bits = make([]uint64, (bv.size+(64-1))/64)
	return io.ReadFull(r, uint64SliceAsByteSlice(bv.bits))
}

// borrowed from roaring
func uint64SliceAsByteSlice(slice []uint64) []byte {
	// make a new slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&slice))

	// update its capacity and length
	header.Len *= 8
	header.Cap *= 8

	// return it
	return *(*[]byte)(unsafe.Pointer(&header))
}

func byteSliceAsUint64Slice(b []byte) []uint64 {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))

	header.Len /= 8
	header.Cap /= 8

	return *(*[]uint64)(unsafe.Pointer(&header))
}

var masks = [65]uint64{
	0,
	(1 << 1) - 1,
	(1 << 2) - 1,
	(1 << 3) - 1,
	(1 << 4) - 1,
	(1 << 5) - 1,
	(1 << 6) - 1,
	(1 << 7) - 1,
	(1 << 8) - 1,
	(1 << 9) - 1,
	(1 << 10) - 1,
	(1 << 11) - 1,
	(1 << 12) - 1,
	(1 << 13) - 1,
	(1 << 14) - 1,
	(1 << 15) - 1,
	(1 << 16) - 1,
	(1 << 17) - 1,
	(1 << 18) - 1,
	(1 << 19) - 1,
	(1 << 20) - 1,
	(1 << 21) - 1,
	(1 << 22) - 1,
	(1 << 23) - 1,
	(1 << 24) - 1,
	(1 << 25) - 1,
	(1 << 26) - 1,
	(1 << 27) - 1,
	(1 << 28) - 1,
	(1 << 29) - 1,
	(1 << 30) - 1,
	(1 << 31) - 1,
	(1 << 32) - 1,
	(1 << 33) - 1,
	(1 << 34) - 1,
	(1 << 35) - 1,
	(1 << 36) - 1,
	(1 << 37) - 1,
	(1 << 38) - 1,
	(1 << 39) - 1,
	(1 << 40) - 1,
	(1 << 41) - 1,
	(1 << 42) - 1,
	(1 << 43) - 1,
	(1 << 44) - 1,
	(1 << 45) - 1,
	(1 << 46) - 1,
	(1 << 47) - 1,
	(1 << 48) - 1,
	(1 << 49) - 1,
	(1 << 50) - 1,
	(1 << 51) - 1,
	(1 << 52) - 1,
	(1 << 53) - 1,
	(1 << 54) - 1,
	(1 << 55) - 1,
	(1 << 56) - 1,
	(1 << 57) - 1,
	(1 << 58) - 1,
	(1 << 59) - 1,
	(1 << 60) - 1,
	(1 << 61) - 1,
	(1 << 62) - 1,
	(1 << 63) - 1,
	(1 << 64) - 1,
}
