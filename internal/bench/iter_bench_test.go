package bench_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkCreate(b *testing.B) {
	b.ReportAllocs()

	var f pi.Iterator
	for i := 0; i < b.N; i++ {
		f = stopImmediately() // 0 allocs, 1ns
	}
	_ = f
}

func BenchmarkCreate2(b *testing.B) {
	b.ReportAllocs()

	var f pi.Iterator
	for i := 0; i < b.N; i++ {
		f = stopAfterOneYield() // 2 allocs, 38ns (still 100x faster than coro.Routine)
	}
	_ = f
}

func BenchmarkCreate3(b *testing.B) {
	b.ReportAllocs()

	var f pi.Iterator
	for i := 0; i < b.N; i++ {
		obj := &coroutineObject{}
		f = obj.Resume // 2 allocs, 38ns
	}
	_ = f
}

func BenchmarkResume(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := stopAfterOneYield() // 2 allocs, 38ns
		r()                      // 0 allocs, 2ns
	}
}

func BenchmarkIteratorsAppend(b *testing.B) {
	b.ReportAllocs()
	var iterators pi.Iterators
	for i := 0; i < 1024; i++ {
		iterators = append(iterators, stopImmediately())
	}
	iterators = iterators.Next()

	b.ResetTimer()

	for i := 0; i < b.N; i += 1024 {
		for j := 0; j < 1024; j++ {
			iterators = append(iterators, stopAfterOneYield())
		}
		iterators = iterators.Next()
	}

	_ = iterators
}

//go:noinline
func stopImmediately() pi.Iterator {
	return func() bool {
		return false
	}
}

//go:noinline
func stopAfterOneYield() pi.Iterator {
	i := 0
	return func() bool {
		i++
		return i <= 1
	}
}

type coroutineObject struct {
	i int
}

func (c *coroutineObject) Resume() bool {
	c.i++
	return c.i <= 1
}
