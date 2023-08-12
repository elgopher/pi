// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package image provides API for decoding images.
//
// Package is used internally by Pi, but it can also be used when writing unit tests.
package image

import (
	"fmt"

	"github.com/elgopher/pi/internal/sfmt"
)

// Image contains information about decoded image.
type Image struct {
	Width, Height int
	// Palette array is filled with black color (#000000)
	// if file has fewer colors than 256.
	Palette [256]RGB
	// Each pixel is a color from 0 to 255.
	// 0th element of slice represent pixel color in top-left corner.
	// 1st element is a next pixel on the right and so on.
	Pixels []byte
}

// String returns Image as string for debugging purposes.
func (i Image) String() string {
	return fmt.Sprintf("{width:%d, height:%d, palette: %+v, pixels:%s}",
		i.Width, i.Height, sfmt.FormatBigSlice(i.Palette[:], 32), sfmt.FormatBigSlice(i.Pixels, 1000))
}

// RGB represents color
type RGB struct{ R, G, B byte }

func (r RGB) String() string {
	out := make([]byte, 7)
	out[0] = '#'
	writeByteAsHex(r.R, out[1:3])
	writeByteAsHex(r.G, out[3:5])
	writeByteAsHex(r.B, out[5:7])
	return string(out)
}

func writeByteAsHex(number byte, out []byte) {
	out[0] = ascii(number / 16)
	out[1] = ascii(number % 16)
}

func ascii(digit byte) byte {
	const asciiA, ascii0 = 65, 48
	if digit > 9 {
		return asciiA + digit - 10
	}
	return ascii0 + digit
}
