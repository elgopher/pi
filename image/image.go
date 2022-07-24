// Package image provides API for decoding images.
//
// Package is used internally by Pi, but it can also be used when writing unit tests.
package image

import "fmt"

// RGB represents color
type RGB struct{ R, G, B byte }

func (r RGB) String() string {
	var rgb = int(r.R)<<16 + int(r.G)<<8 + int(r.B)
	return fmt.Sprintf("#%.6x", rgb)
}

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
