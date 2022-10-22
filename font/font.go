// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package font

import (
	"bytes"
	"fmt"

	"github.com/elgopher/pi/image"
)

// Load loads font-sheet (png image) and converts
// it to font data. Image must be 128x128. Each char is 8x8.
// Char 0 is in the top-left corner. Char 1 to the right.
//
// This function can be used if you want to use 3rd font:
//
// myFont := pi.Font{Width:4, WidthSpecial:8, Height: 8}
// pi.Load(png, myFont.Data[:])
//
// Color with index 0 is treated as background. Any other color
// as foreground.
//
// The result is inserted into fontData slice. The size of slice must
// be 2048.
func Load(png []byte, fontData []byte) error {
	img, err := image.DecodePNG(bytes.NewReader(png))
	if err != nil {
		return fmt.Errorf("decoding font failed: %w", err)
	}

	if err = load(img, fontData[:]); err != nil {
		return fmt.Errorf("error system font: %w", err)
	}

	return nil
}

func load(img image.Image, fontData []byte) error {
	const (
		charWidth, charHeight  = 8, 8   // in pixels
		rows, cells            = 16, 16 // number of rows and cells in font sheet
		imgWidth, imgHeight    = charWidth * cells, charHeight * rows
		expectedNumberOfPixels = imgWidth * imgHeight
		charBytes              = 8 // how many bytes a single char occupies in fontData
		numberOfChars          = 256
		expectedFontDataLen    = numberOfChars * charBytes
	)

	if img.Width != imgWidth || img.Height != imgHeight {
		return fmt.Errorf("invalid font image size: must be %dx%d", imgWidth, imgHeight)
	}
	if len(img.Pixels) != expectedNumberOfPixels {
		return fmt.Errorf("invalid font image pixels slice len: must be %d", expectedNumberOfPixels)
	}
	if len(fontData) != expectedFontDataLen {
		return fmt.Errorf("invalid fontData len: must be %d", expectedFontDataLen)
	}

	clear(fontData)

	for row := 0; row < rows; row++ {
		for cell := 0; cell < cells; cell++ {

			for y := 0; y < charHeight; y++ {
				for x := 0; x < charWidth; x++ {
					imageOffsetY := row*imgWidth*charHeight + y*imgWidth
					imageOffset := imageOffsetY + (cell * charWidth) + x

					if img.Pixels[imageOffset] != 0 {
						outOffset := row*charBytes*cells + y + cell*charBytes
						fontData[outOffset] |= 1 << x
					}
				}
			}
		}
	}

	return nil
}

func clear(data []byte) {
	zeroData := make([]byte, len(data))
	copy(data, zeroData)
}
