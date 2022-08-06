// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package font

import (
	"fmt"

	"github.com/elgopher/pi/image"
)

// Load gets non-zero pixels from the font sheet and converts them to font data understood by π.
// The result is inserted into fontData slice.
//
// Img must be 128x128. Each char is 8x8. Char 0 is in the top-left corner. Char 1 to the right.
//
// fontData must have length of 2048.
func Load(img image.Image, fontData []byte) error {
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
