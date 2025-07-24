// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pifont_test

import (
	_ "embed"
	"github.com/elgopher/pi/pitest"
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pifont"
)

//go:embed internal/test/font.png
var fontPNG []byte

var fontSheet pifont.Sheet

func init() {
	prevPalette := pi.Palette
	defer func() {
		pi.Palette = prevPalette
	}()

	pi.Palette = pi.DecodePalette(fontPNG)
	fontCanvas := pi.DecodeCanvas(fontPNG)
	fontSheet = pifont.Sheet{
		Chars: map[rune]pi.Sprite{
			'S': {
				Area:   pi.Area[int]{X: 0, Y: 0, W: 8, H: 8},
				Source: fontCanvas,
			},
			'T': {
				Area:   pi.Area[int]{X: 8, Y: 0, W: 8, H: 8},
				Source: fontCanvas,
			},
			'⬤': {
				Area:   pi.Area[int]{X: 0, Y: 8, W: 8, H: 8},
				Source: fontCanvas,
			},
			'❤': {
				Area:   pi.Area[int]{X: 8, Y: 8, W: 8, H: 8},
				Source: fontCanvas,
			},
		},
		Height:  8,
		FgColor: 1,
		BgColor: 0,
	}
}

var (
	//go:embed internal/test/text-color-equal-to-bg.png
	textColorEqualToBg []byte
	//go:embed internal/test/text-color-equal-to-fg.png
	textColorEqualToFg []byte
	//go:embed internal/test/text-color-different.png
	textColorDifferent []byte
)

func TestSheet_Print(t *testing.T) {
	t.Run("should print with different colors", func(t *testing.T) {
		tests := map[string]struct {
			bgColor   pi.Color
			textColor pi.Color
			png       []byte
		}{
			"text color different than Bg and Fg": {
				bgColor:   0,
				textColor: 2,
				png:       textColorDifferent,
			},
			"text color equal to Bg": {
				bgColor:   2,
				textColor: 0,
				png:       textColorEqualToBg,
			},
			"text color equal to Fg": {
				bgColor:   0,
				textColor: 1,
				png:       textColorEqualToFg,
			},
		}

		for testName, testCase := range tests {
			t.Run(testName, func(t *testing.T) {
				pi.Palette = pi.DecodePalette(testCase.png)
				expectedCanvas := pi.DecodeCanvas(testCase.png)
				pi.SetScreenSize(8, 8)

				pi.Screen().Clear(testCase.bgColor)
				pi.SetColor(testCase.textColor)
				pi.SetTransparency(testCase.textColor, false)
				// when
				fontSheet.Print("S", 0, 0)
				// then
				pitest.AssertSurfaceEqual(t, expectedCanvas, pi.Screen())
			})
		}
	})
}

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
