// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkSset(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 1000; i++ { // Sset is too fast
			pi.Sset(2, 2, 7)
		}
	})
}

func BenchmarkSget(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 1000; i++ { // Sget is too fast
			sink = pi.Sget(2, 2)
		}
	})
}
