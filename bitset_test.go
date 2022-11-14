package bitset

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestSize(t *testing.T) {
	b1, _ := New(7)
	if b1.BitsCount() != 7 {
		t.Errorf("should have %v bits, got %v", 7, b1.BitsCount())
	}
	if b1.BytesCount() != 1 {
		t.Errorf("should have %v bytes, got %v", 1, b1.BytesCount())
	}

	b2, _ := New(8)
	if b2.BitsCount() != 8 {
		t.Errorf("should have %v bits, got %v", 8, b2.BitsCount())
	}
	if b2.BytesCount() != 1 {
		t.Errorf("should have %v bytes, got %v", 1, b2.BytesCount())
	}

	b3, _ := New(9)
	if b3.BitsCount() != 9 {
		t.Errorf("should have %v bits, got %v", 9, b3.BitsCount())
	}
	if b3.BytesCount() != 2 {
		t.Errorf("should have %v bytes, got %v", 2, b3.BytesCount())
	}
}

func TestNewFromBytes(t *testing.T) {
	b1, _ := New(1000)
	b1.Set(1)
	b1.Set(100)
	b1.Set(354)
	b2 := NewFromBytes(1000, b1.Bytes())
	if !bytes.Equal(b1.Bytes(), b2.Bytes()) {
		t.Errorf("bytes should be identical")
	}
}

func TestSetAndTest(t *testing.T) {
	b1, _ := New(1000)
	if b1.BytesCount() != 125 {
		t.Errorf("should have %v bytes, got %v", 125, b1.BytesCount())
	}
	if b1.OnesCount() != 0 {
		t.Errorf("should have %v ones, got %v", 0, b1.OnesCount())
	}

	b1.Set(7)
	if b1.OnesCount() != 1 {
		t.Errorf("should have %v ones, got %v", 1, b1.OnesCount())
	}
	if !b1.Test(7) {
		t.Error("should be true, got false")
	}

	b1.Set(27)
	if b1.OnesCount() != 2 {
		t.Errorf("should have %v ones, got %v", 2, b1.OnesCount())
	}
	if !b1.Test(27) {
		t.Error("should be true, got false")
	}
	if b1.Test(230) {
		t.Error("should be false, got true")
	}
}

func TestHexEncode(t *testing.T) {
	b, _ := New(1)
	b.Set(0)
	if b.HexEncode() != "01" {
		t.Errorf("expected 01, got %v", b.HexEncode())
	}
	b.Set(7)
	if b.HexEncode() != "81" {
		t.Errorf("expected 01, got %v", b.HexEncode())
	}
}

func TestUnion(t *testing.T) {
	b1, _ := New(20)
	b1.Set(5)
	b2, _ := New(20)
	b2.Set(6)
	b1.Union(b2)
	if !b1.Test(5) || !b1.Test(6) {
		t.Errorf("expected bits 5 and 6 to be set after Union. b2 = %b, b1 = %b", b2.bytes, b1.bytes)
	}
}

func TestIntersect(t *testing.T) {
	b1, _ := New(20)
	b1.Set(5)
	b2, _ := New(20)
	b2.Set(5)
	b2.Set(6)
	b1.Intersect(b2)
	if !b1.Test(5) || b1.Test(6) {
		t.Errorf("expected bit 5 and not bit 6 to be set after Intersect. b2 = %b, b1 = %b", b2.bytes, b1.bytes)
	}
}

func FuzzUnion(f *testing.F) {
	testcases := []uint64{10, 20, 30, 40, 50}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, size uint64) {
		t.Logf("size = %v", size)

		b1, err := New(size)
		if size <= 0 && err == nil {
			t.Errorf("should return error for negative size, size = %v", size)
		}

		b2, err := New(uint64(size))
		if size <= 0 && err == nil {
			t.Errorf("should return error for negative size, size = %v", size)
		}

		if size <= 0 && err != nil {
			t.Skip()
		}

		bit1 := uint64(rand.Int63n(int64(size)))
		bit2 := uint64(rand.Int63n(int64(size)))
		if bit2 == bit1 {
			t.Skip()
		}

		b1.Set(bit1)
		b2.Set(bit2)
		b1.Union(b2)
		if !(b1.Test(bit1) && b1.Test(bit2)) {
			t.Error("failed test")
		}
		if b1.OnesCount() != 2 {
			t.Error("Should only have 2 bits set to 1")
		}

	})
}

func FuzzIntersect(f *testing.F) {
	testcases := []uint64{10, 20, 30, 40, 50}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, size uint64) {
		b1, err := New(size)
		if size <= 0 && err == nil {
			t.Errorf("should return error for negative size, size = %v", size)
		}

		b2, err := New(uint64(size))
		if size <= 0 && err == nil {
			t.Errorf("should return error for negative size, size = %v", size)
		}

		if size <= 0 && err != nil {
			t.Skip()
		}

		bit1 := uint64(rand.Int63n(int64(size)))
		bit2 := uint64(rand.Int63n(int64(size)))
		bit3 := uint64(rand.Int63n(int64(size)))
		if bit1 == bit2 || bit1 == bit3 {
			t.Skip()
		}

		b1.Set(bit1)
		b1.Set(bit2)
		b2.Set(bit2)
		b2.Set(bit3)
		b1.Intersect(b2)
		if !b2.Test(bit2) {
			t.Error("failed test")
		}
		if b1.OnesCount() != 1 {
			t.Error("Should only have 1 bit set to 1")
		}

	})
}
