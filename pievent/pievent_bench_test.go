// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pievent_test

import (
	"testing"

	"github.com/elgopher/pi/pievent"
)

func init() {
	pievent.GlobalTracingOff = true
}

func BenchmarkPublish(b *testing.B) {
	b.ReportAllocs()
	target := pievent.NewTarget[someEvent]()
	event := someEvent{a: "1"} // event should not be empty for get meaningful results
	target.SubscribeAll(func(someEvent, pievent.Handler) {})

	for b.Loop() {
		target.Publish(event) // zero alokacji! LOVE IT
	}
}

func BenchmarkSubscribe(b *testing.B) {
	b.ReportAllocs()
	target := pievent.NewTarget[someEvent]()
	listener := func(someEvent, pievent.Handler) {}

	for b.Loop() {
		// 3 allocs, because stack trace is analyzed - only for debugging
		// 0 allocs for production code
		target.SubscribeAll(listener)
	}
}

func BenchmarkSubscribeEvent(b *testing.B) {
	b.ReportAllocs()
	target := pievent.NewTarget[someEvent]()
	listener := func(someEvent, pievent.Handler) {}

	for b.Loop() {
		// 3 allocs, because stack trace is analyzed - only for debugging
		// 0 allocs for production code
		target.Subscribe(someEvent{a: "a"}, listener)
	}
}

type someEvent struct {
	a string
}
