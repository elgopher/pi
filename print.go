// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/elgopher/pi/font"
	"github.com/elgopher/pi/mem"
)

// Print prints text on the screen using system font. It takes into consideration
// clipping region and camera position.
//
// Only unicode characters with code < 256 are supported. Unsupported chars
// are printed as question mark. The entire table of available chars can be
// found here: https://github.com/elgopher/pi/blob/master/internal/system-font.png
//
// Print returns the right-most x position that occurred while printing.
func Print(text string, x, y int, color byte) (rightMostX int) {
	return Font(mem.SystemFont).Print(text, x, y, color)
}

// PrintCustom prints text in the same way as Print, but using custom font.
func PrintCustom(text string, x, y int, color byte) (rightMostX int) {
	// FIXME Probably escape character should be used to switch the font instead
	return Font(mem.CustomFont).Print(text, x, y, color)
}

//go:embed internal/system-font.png
var systemFontPNG []byte

// Font contains all information about loaded font and provides method to Print the text.
type Font mem.Font

// Print prints text on the screen at given coordinates. It takes into account
// clipping region and camera position.
//
// Only unicode characters with code < 256 are supported. Unsupported chars
// are printed as question mark. The entire table of available chars can be
// found here: https://github.com/elgopher/pi/blob/master/internal/system-font.png
//
// Print returns the right-most x position that occurred while printing.
func (f Font) Print(text string, x, y int, color byte) int {
	startX := x

	for _, r := range text {
		if r != '\n' {
			width := f.printRune(r, x, y, color)
			x += width
		} else {
			x = startX
			y += f.Height
		}
	}

	return x
}

func (f Font) printRune(r rune, sx, sy int, color byte) int {
	if r > 255 {
		r = '?'
	}

	index := int(r) * 8

	for y := 0; y < 8; y++ {
		if mem.ClippingRegion.Y > sy+y-mem.Camera.Y {
			continue
		}
		if mem.ClippingRegion.Y+mem.ClippingRegion.H <= sy+y-mem.Camera.Y {
			continue
		}

		offset := mem.ScreenWidth*y + sx + sy*mem.ScreenWidth - mem.Camera.Y*mem.ScreenWidth - mem.Camera.X
		line := f.Data[index+y]
		for bit := 0; bit < 8; bit++ {
			if mem.ClippingRegion.X > sx+bit-mem.Camera.X {
				continue
			}
			if mem.ClippingRegion.X+mem.ClippingRegion.W <= sx+bit-mem.Camera.X {
				continue
			}
			if line&(1<<bit) == 1<<bit {
				mem.ScreenData[offset+bit] = color
			}
		}
	}

	if r < 128 {
		return f.Width
	} else {
		return f.WidthSpecial
	}
}

func loadCustomFont(resources fs.ReadFileFS) error {
	fileContents, err := resources.ReadFile("custom-font.png")
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error loading custom-font.png file: %w", err)
	}

	return font.Load(fileContents, mem.CustomFont.Data[:])
}
