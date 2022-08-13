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
	Height       int
}

func loadSystemFont() error {
	img, err := image.DecodePNG(bytes.NewReader(systemFontPNG))
	if err != nil {
		return fmt.Errorf("decoding system font failed: %w", err)
	}

	if err = font.Load(img, systemFont.Data[:]); err != nil {
		return fmt.Errorf("error loading system font: %w", err)
	}

	return nil
}

var cursor pos

// Cursor set cursor position (in pixels) used by Print.
//
// Cursor returns previously set cursor position.
func Cursor(x, y int) (prevX, prevY int) {
	prevX, prevY = cursor.x, cursor.y
	cursor.x = x
	cursor.y = y
	return
}

// CursorReset resets cursor position used by Print to 0,0.
//
// CursorReset returns previously set cursor position.
func CursorReset() (prevX, prevY int) {
	return Cursor(0, 0)
}

// Print prints text on the screen. It takes into consideration cursor position,
// color, clipping region and camera position.
//
// After printing all characters Print goes to the next line. When there is no space
// left on screen the clipping region is scrolled to make room.
//
// Only unicode characters with code < 256 are supported. Unsupported chars
// are printed as question mark. The entire table of available chars can be
// found here: https://github.com/elgopher/pi/blob/master/internal/system-font.png
//
// Print returns the right-most x position that occurred while printing.
func Print(text string) (x int) {
	if cursor.y > scrHeight-systemFont.Height {
		lines := systemFont.Height - (scrHeight - cursor.y)
		scroll(lines)
	}

	startingX := cursor.x
	for _, r := range text {
		printRune(r)
	}

	x = cursor.x
	cursor.x = startingX
	cursor.y += systemFont.Height
	if cursor.y > scrHeight-systemFont.Height {
		scroll(systemFont.Height)
	}

	return
}

func scroll(lines int) {
	if scrHeight <= lines {
		Cls()
		cursor.y = scrHeight - systemFont.Height
		return
	}

	for y := clippingRegion.y; y < scrHeight-lines; y++ {
		srcOffset := y*scrWidth + clippingRegion.x
		dstOffset := srcOffset + lines*scrWidth
		copy(ScreenData[srcOffset:], ScreenData[dstOffset:dstOffset+clippingRegion.w])
	}

	for y := scrHeight - lines; y < clippingRegion.y+clippingRegion.h; y++ {
		offset := y*scrWidth + clippingRegion.x
		copy(ScreenData[offset:], zeroScreenData[:clippingRegion.w])
	}

	cursor.y -= lines
}

func printRune(r rune) {
	if r > 255 {
		r = '?'
	}

	index := int(r) * 8

	for y := 0; y < 8; y++ {
		if clippingRegion.y > cursor.y+y-camera.y {
			continue
		}
		if clippingRegion.y+clippingRegion.h <= cursor.y+y-camera.y {
			continue
		}

		offset := scrWidth*y + cursor.x + cursor.y*scrWidth - camera.y*scrWidth - camera.x
		line := systemFont.Data[index+y]
		for bit := 0; bit < 8; bit++ {
			if clippingRegion.x > cursor.x+bit-camera.x {
				continue
			}
			if clippingRegion.x+clippingRegion.w <= cursor.x+bit-camera.x {
				continue
			}
			if line&(1<<bit) == 1<<bit {
				ScreenData[offset+bit] = color
			}
		}
	}

	if r < 128 {
		cursor.x += systemFont.Width
	} else {
		cursor.x += systemFont.WidthSpecial
	}
}
