// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package bench_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkMinInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 60; j++ {
			pi.MinInt(j, j+1)
		}
		for j := 0; j < 60; j++ {
			pi.MinInt(j+1, j)
		}
	}
}

func BenchmarkMaxInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 60; j++ {
			pi.MaxInt(j, j+1)
		}
		for j := 0; j < 60; j++ {
			pi.MaxInt(j+1, j)
		}
	}
}

func BenchmarkMidInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			pi.MidInt(j, j+1, j+2) // x
			pi.MidInt(j+2, j+1, j) // x
			pi.MidInt(j+1, j, j+2) // y
			pi.MidInt(j+1, j+2, j) // y
			pi.MidInt(j, j+2, j+1) // z
			pi.MidInt(j+2, j, j+1) // z
		}
	}
}
