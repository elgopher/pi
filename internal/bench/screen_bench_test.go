// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package bench_test

import (
	"fmt"
	"testing"

	"github.com/elgopher/pi"
)

const color = 7

func BenchmarkCls(b *testing.B) {
	runBenchmarks(b, func(r Resolution) {
		pi.Cls()
	})
}

func BenchmarkClsCol(b *testing.B) {
	runBenchmarks(b, func(r Resolution) {
		pi.ClsCol(color)
	})
}

func BenchmarkPset(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 1000; i++ { // Pset is too fast
			pi.Pset(2, 2, color)
		}
	})
}

var sink byte

func BenchmarkPget(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 1000; i++ { // Pget is too fast
			sink = pi.Pget(2, 2)
		}
	})
}

func BenchmarkSpr(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ { // Spr is too fast
			pi.Spr(0, 16, 16)
		}
	})
}

func BenchmarkSprSize(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ { // SprSize is too fast
			pi.SprSize(0, 16, 16, 2.0, 2.0)
		}
	})
}

func BenchmarkSprSizeFlip(b *testing.B) {
	b.Run("flip x only", func(b *testing.B) {
		runBenchmarks(b, func(res Resolution) {
			for i := 0; i < 100; i++ { // SprSizeFlip is too fast
				pi.SprSizeFlip(0, 16, 16, 2.0, 2.0, true, false)
			}
		})
	})
	b.Run("flip y only", func(b *testing.B) {
		runBenchmarks(b, func(res Resolution) {
			for i := 0; i < 100; i++ { // SprSizeFlip is too fast
				pi.SprSizeFlip(0, 16, 16, 2.0, 2.0, false, true)
			}
		})
	})
	b.Run("flip xy", func(b *testing.B) {
		runBenchmarks(b, func(res Resolution) {
			for i := 0; i < 100; i++ { // SprSizeFlip is too fast
				pi.SprSizeFlip(0, 16, 16, 2.0, 2.0, true, true)
			}
		})
	})
}

func BenchmarkPalt(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ { // Palt is too fast
			pi.Palt(byte(i), true)
		}
	})
}

func BenchmarkPaltReset(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ { // PaltReset is too fast
			pi.PaltReset()
		}
	})
}

func BenchmarkPal(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ { // Pal is too fast
			pi.Pal(byte(i), byte(i+1))
		}
	})
}

func BenchmarkPalDisplay(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ { // PalDisplay is too fast
			pi.PalDisplay(byte(i), byte(i+1))
		}
	})
}

func BenchmarkPalReset(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ { // PalReset is too fast
			pi.PalReset()
		}
	})
}

type Resolution struct{ W, H int }

func (s Resolution) String() string {
	return fmt.Sprintf("%dx%d", s.W, s.H)
}

func runBenchmarks(b *testing.B, callback func(res Resolution)) {
	var resolutions = [...]Resolution{
		{W: 128, H: 128},
		{W: 256, H: 256},
		{W: 512, H: 512},
	}

	for _, resolution := range resolutions {
		b.Run(resolution.String(), func(b *testing.B) {
			pi.Reset()
			pi.SetScreenSize(resolution.W, resolution.H)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				callback(resolution)
			}
		})
	}
}
