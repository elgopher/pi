// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pifont_test

import (
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pifont"
)

func BenchmarkSheet_Print(b *testing.B) {
	sheet := pifont.Sheet{
		Chars: map[rune]pi.Sprite{
			'a': {
				Area:   pi.Area[int]{X: 0, Y: 0, W: 8, H: 8},
				Source: pi.NewCanvas(8, 8),
			},
		},
	}
	b.ReportAllocs()
	for b.Loop() {
		sheet.Print("aaaaaaaaaaaaaaaaaaaa", 100, 100)
	}
}
