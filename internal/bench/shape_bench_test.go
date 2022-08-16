// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package bench_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkRectFill(b *testing.B) {
	runBenchmarks(b, func(r Resolution) {
		pi.RectFill(0, 0, r.W-1, r.H-1, color)
	})
}

func BenchmarkRect(b *testing.B) {
	runBenchmarks(b, func(r Resolution) {
		pi.Rect(0, 0, r.W-1, r.H-1, color)
	})
}

func BenchmarkLine(b *testing.B) {
	b.Run("slope=1", func(b *testing.B) {
		runBenchmarks(b, func(r Resolution) {
			pi.Line(0, 0, r.W-1, r.H-1, color)
		})
	})
	b.Run("slope>1", func(b *testing.B) {
		runBenchmarks(b, func(r Resolution) {
			pi.Line(0, 0, r.W-2, r.H-1, color)
		})
	})
	b.Run("vertical", func(b *testing.B) {
		runBenchmarks(b, func(r Resolution) {
			pi.Line(0, 0, 0, r.H-1, color)
		})
	})
	b.Run("horizontal", func(b *testing.B) {
		runBenchmarks(b, func(r Resolution) {
			pi.Line(0, 0, r.W-1, 0, color)
		})
	})
}

func BenchmarkCirc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runBenchmarks(b, func(r Resolution) {
			pi.Circ(r.W/2, r.H/2, r.W/2, color)
		})
	}
}

func BenchmarkCircFill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runBenchmarks(b, func(r Resolution) {
			pi.CircFill(r.W/2, r.H/2, r.W/2, color)
		})
	}
}
