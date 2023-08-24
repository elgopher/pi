package coro_test

import (
	"testing"

	"github.com/elgopher/pi/coro"
)

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	var r *coro.Routine[struct{}]

	for i := 0; i < b.N; i++ {
		r = coro.New(f2) // 7 allocs :( 4us on windows :( But on linux it is 1us and 5 allocs!
	}

	_ = r
}

func BenchmarkCreate(b *testing.B) {
	b.ReportAllocs()

	var r *coro.Routine[struct{}]

	for i := 0; i < b.N; i++ {
		r = coro.WithReturn(f) // 6 allocs :( 4us on windows :( But on linux it is 1us and 5 allocs!
	}

	_ = r
}

func BenchmarkResume(b *testing.B) {
	b.ReportAllocs()

	var r *coro.Routine[struct{}]

	for i := 0; i < b.N; i++ {
		r = coro.WithReturn(f) // 6 allocs
		r.Resume()             // 1 alloc, 0.8us :(
	}
	_ = r
}

func BenchmarkResumeUntilFinish(b *testing.B) {
	b.ReportAllocs()

	var r *coro.Routine[struct{}]

	for i := 0; i < b.N; i++ {
		r = coro.WithReturn(f) // 6 allocs
		r.Resume()             // 1 alloc, 0.8us :(
		r.Resume()             // 1 alloc, 0.8us :(
	}
	_ = r
}

func BenchmarkCancel(b *testing.B) {
	b.ReportAllocs()

	var r *coro.Routine[struct{}]

	for i := 0; i < b.N; i++ {
		r = coro.WithReturn(f) // 6 allocs
		r.Cancel()             // -2 alloc????
	}
	_ = r
}

//go:noinline
func f2(yield coro.Yield) {
	yield()
}

//go:noinline
func f(yield coro.YieldReturn[struct{}]) {
	yield(struct{}{})
}
