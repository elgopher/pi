// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package bench_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkRectFill(b *testing.B) {
	runBenchmarks(b, func(r Resolution) {
		pi.RectFill(0, 0, r.W-1, r.H-1)
	})
}

func BenchmarkRect(b *testing.B) {
	runBenchmarks(b, func(r Resolution) {
		pi.Rect(0, 0, r.W-1, r.H-1)
	})
}
