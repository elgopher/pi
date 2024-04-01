// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package bench_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkPixMapPointer(b *testing.B) {
	pixMap := pi.NewPixMap(3, 2) // 3x2

	var ptr pi.Pointer
	var ok bool

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ptr, ok = pixMap.Pointer(1, 1, 2, 1)
	}
	_ = ptr
	_ = ok
}

func BenchmarkCopy(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ {
			pi.SprSheet().WithClip(0, 0, 16, 16).Copy(pi.Scr(), 16, 16) // 2x times faster than SprSize
		}
	})
}

func SrcAtop(dst, src []byte) { copy(dst, src) }

func BenchmarkMerge(b *testing.B) {
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ {
			pi.SprSheet().WithClip(0, 0, 16, 16).Merge(pi.Scr(), 16, 16, SrcAtop) // FAST! NOT AS FAST AS COPY BUT THE PERF IS GREAT!
		}
	})
}

func BenchmarkForeach(b *testing.B) {
	src := make([]byte, 16)
	runBenchmarks(b, func(res Resolution) {
		for i := 0; i < 100; i++ {
			pi.Scr().WithClip(0, 0, 16, 16).Foreach(func(x, y int, dst []byte) { copy(dst, src) })
		}
	})
}
