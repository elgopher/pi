// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package bench_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkMidInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			pi.Mid(j, j+1, j+2) // y
			pi.Mid(j+2, j+1, j) // y
			pi.Mid(j+1, j, j+2) // x
			pi.Mid(j+1, j+2, j) // x
			pi.Mid(j, j+2, j+1) // z
			pi.Mid(j+2, j, j+1) // z
		}
	}
}

func BenchmarkMid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			f := float64(j)
			pi.Mid(f, f+1, f+2) // y
			pi.Mid(f+2, f+1, f) // y
			pi.Mid(f+1, f, f+2) // x
			pi.Mid(f+1, f+2, f) // x
			pi.Mid(f, f+2, f+1) // z
			pi.Mid(f+2, f, f+1) // z
		}
	}
}
