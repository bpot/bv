package bv

import "testing"

func BenchmarkGetInt(b *testing.B) {
	bv := New(6 * 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := uint(0); j < 1024; j++ {
			bv.GetInt(j*6, 6)
		}
	}
}
