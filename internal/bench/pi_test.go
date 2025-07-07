// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package bench_test

import (
	"math/rand"
	"testing"

	"github.com/elgopher/pi"
)

func BenchmarkSpr(b *testing.B) {
	pi.SetScreenSize(256, 256)
	canvas := pi.NewCanvas(256, 256)
	for i := 0; i < canvas.W()*canvas.H(); i++ {
		canvas.Data()[i] = pi.Color(rand.Intn(256))
	}

	sprite := pi.SpriteFrom(canvas, 128, 128, 16, 16) //  396.6 (1 color table), 510 (4), vs 559 (ReadMask and TargetMask)

	for b.Loop() {
		pi.Spr(sprite, 10, 10)
	}
}

func BenchmarkLine(b *testing.B) {
	pi.SetScreenSize(256, 256)
	for b.Loop() {
		pi.Line(64, 64, 30, 30) // 94
	}
}

func BenchmarkRect(b *testing.B) {
	pi.SetScreenSize(256, 256)
	b.ReportAllocs()
	for b.Loop() {
		pi.RectFill(64, 64, 128, 128) // 2771 (one color table) vs 4337 (4 color tables) vs 4833 (Write and ReadMask)
	}
}
