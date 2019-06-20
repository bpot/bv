package bv

import (
	"bytes"
	"testing"
)

func TestSerialization(t *testing.T) {
	bv := New(1024)
	var buf bytes.Buffer
	_, _ = bv.WriteTo(&buf)

	if 128+8 != len(buf.Bytes()) {
		t.Errorf("expected 136; got %d", len(buf.Bytes()))
	}
}

func TestByteBacked(t *testing.T) {
	v := New(1024)
	v.Set(1, true)
	v.Set(512, true)
	v.Set(1000, true)

	var buf bytes.Buffer
	_, _ = v.WriteTo(&buf)

	unserialized, serializedLength, err := NewByteBacked(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if uint64(len(buf.Bytes())) != serializedLength {
		t.Errorf("expected %d; got %d", len(buf.Bytes()), serializedLength)
	}

	if !unserialized.Get(1) {
		t.Errorf("expected 1 to be set")
	}

	if !unserialized.Get(512) {
		t.Errorf("expected 512 to be set")
	}

	if !unserialized.Get(1000) {
		t.Errorf("expected 1000 to be set")
	}
}
