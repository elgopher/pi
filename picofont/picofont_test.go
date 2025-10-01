// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package picofont_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/pitest"
)

//go:embed "font.png"
var fontPNG []byte

func TestPrint(t *testing.T) {
	t.Run("should print each character", func(t *testing.T) {
		pi.SetScreenSize(128, 128)

		pi.Palette = pi.DecodePalette(fontPNG)
		canvas := pi.DecodeCanvas(fontPNG)

		var table strings.Builder

		// print narrow characters
		for i := 16; i < 128; i++ { // skip escape codes below 16 (such as LF)
			table.WriteRune(rune(i))
			table.WriteByte(' ')
			if i%16 == 15 {
				table.WriteByte('\n')
			}
		}
		// print wide characters
		for i := 128; i < 256; i++ {
			table.WriteRune(rune(i))
			if i%16 == 15 {
				table.WriteByte('\n')
			}
		}
		pi.SetColor(1)
		// when
		picofont.Print(table.String(), 0, 8)
		// then
		pitest.AssertSurfaceEqual(t, canvas, pi.Screen())
	})
}
