// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkPrint(b *testing.B) {
	runBenchmarks(b, func(r Resolution) {
		for j := 0; j < 10; j++ {
			pi.Cursor(0, 0)
			pi.Print("Hello")
		}
	})
}

func BenchmarkPrintWithScroll(b *testing.B) {
	runBenchmarks(b, func(r Resolution) {
		pi.Print("Hello")
	})
}
