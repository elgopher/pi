// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/elgopher/pi/font"
)

const fontDataSize = 8 * 256

var (
	systemFont = Font{
		Width:        4,
		SpecialWidth: 8,
		Height:       6,
	}

	//go:embed internal/system-font.png
	systemFontPNG []byte

	defaultCustomFont = Font{Width: 4, SpecialWidth: 8, Height: 6}

	customFont = Font{
		Data:         make([]byte, fontDataSize),
		Width:        defaultCustomFont.Width,
		SpecialWidth: defaultCustomFont.SpecialWidth,
		Height:       defaultCustomFont.Height,
	}
)

func init() {
	var err error
	systemFont.Data, err = font.Load(systemFontPNG)
	if err != nil {
		panic(err)
	}
}

// Font contains all information about loaded font and  provides method to Print the text.
type Font struct {
	// Data contains all 256 characters sorted by their ascii-like number.
	// Each character is 8 subsequent bytes, starting from the top.
	// Left-most pixel in a line is bit 0. Right-most pixel in a line is bit 7.
	//
	// The size of slice is always 8 * 256 = 2048.
	//
	// Can be freely read and updated. Changes will be visible immediately.
	Data []byte
	// Width in pixels for all characters below 128. For Width > 8 only 8 pixels are drawn.
	Width int
	// SpecialWidth is a with of all special characters (code>=128)
	// For SpecialWidth > 8 only 8 pixels are drawn.
	SpecialWidth int
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

	clip := screen.clip

	for y := 0; y < 8; y++ {
		if clip.Y > sy+y-ScreenCamera.Y {
			continue
		}
		if clip.Y+clip.H <= sy+y-ScreenCamera.Y {
			continue
		}

		screenWidth := screen.Width()

		offset := screenWidth*y + sx + sy*screenWidth - ScreenCamera.Y*screenWidth - ScreenCamera.X
		line := f.Data[index+y]
		for bit := 0; bit < 8; bit++ {
			if clip.X > sx+bit-ScreenCamera.X {
				continue
			}
			if clip.X+clip.W <= sx+bit-ScreenCamera.X {
				continue
			}
			if line&(1<<bit) == 1<<bit {
				screen.pix[offset+bit] = color
			}
		}
	}

	if r < 128 {
		return f.Width
	} else {
		return f.SpecialWidth
	}
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

func loadCustomFont(resources fs.ReadFileFS) error {
	fileContents, err := resources.ReadFile("custom-font.png")
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error loading custom-font.png file: %w", err)
	}

	customFont.Data, err = font.Load(fileContents)
	return err
}

func SystemFont() Font {
	return systemFont
}

func CustomFont() Font {
	return customFont
}

func SetCustomFontWidth(w int) {
	if w < 0 {
		w = 0
	}
	if w > 8 {
		w = 8
	}
	customFont.Width = w
}

func SetCustomFontSpecialWidth(w int) {
	if w < 0 {
		w = 0
	}
	if w > 8 {
		w = 8
	}
	customFont.SpecialWidth = w
}

func SetCustomFontHeight(height int) {
	if height < 0 {
		height = 0
	}
	customFont.Height = height
}
