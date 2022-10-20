// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/elgopher/pi/image"
	"github.com/elgopher/pi/internal/font"
)

var systemFont = Font{
	Width:        4,
	WidthSpecial: 8,
	Height:       6,
}

//go:embed internal/system-font.png
var systemFontPNG []byte

// Font contains all information about loaded font.
type Font struct {
	// Data contains all 256 characters sorted by their ascii-like number.
	// Each character is 8 subsequent bytes, starting from the top.
	// Left-most pixel in a line is bit 0. Right-most pixel in a line is bit 7.
	Data [8 * 256]byte
	// Width in pixels for all characters below 128
	Width int
	// WidthSpecial is a with of all special characters (code>=128)
	WidthSpecial int
	// Height of line
	Height int
}

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
		if clippingRegion.y > sy+y-camera.y {
			continue
		}
		if clippingRegion.y+clippingRegion.h <= sy+y-camera.y {
			continue
		}

		offset := scrWidth*y + sx + sy*scrWidth - camera.y*scrWidth - camera.x
		line := f.Data[index+y]
		for bit := 0; bit < 8; bit++ {
			if clippingRegion.x > sx+bit-camera.x {
				continue
			}
			if clippingRegion.x+clippingRegion.w <= sx+bit-camera.x {
				continue
			}
			if line&(1<<bit) == 1<<bit {
				ScreenData[offset+bit] = color
			}
		}
	}

	if r < 128 {
		return f.Width
	} else {
		return f.WidthSpecial
	}
}

// loadFontData loads font-sheet (png image) and converts
// it to font data. Image must be 128x128. Each char is 8x8.
// Char 0 is in the top-left corner. Char 1 to the right.
//
// This function can be used if you want to use 3rd font:
//
// myFont := pi.Font{Width:4, WidthSpecial:8,Height: 8}
// pi.loadFontData(png, myFont.Data[:])
func loadFontData(png []byte, out []byte) error {
	img, err := image.DecodePNG(bytes.NewReader(png))
	if err != nil {
		return fmt.Errorf("decoding font failed: %w", err)
	}

	if err = font.Load(img, out[:]); err != nil {
		return fmt.Errorf("error system font: %w", err)
	}

	return nil
}

// Print prints text on the screen using system font. It takes into consideration
// clipping region and camera position.
//
// Only unicode characters with code < 256 are supported. Unsupported chars
// are printed as question mark. The entire table of available chars can be
// found here: https://github.com/elgopher/pi/blob/master/internal/system-font.png
//
// Print returns the right-most x position that occurred while printing.
func Print(text string, x, y int, color byte) (rightMostX int) {
	return systemFont.Print(text, x, y, color)
}
