// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package vm

import (
	"fmt"
)

var (
	// Palette has all colors available in the game. Up to 256.
	// Palette is taken from loaded sprite sheet (which must be
	// a PNG file with indexed color mode). If sprite-sheet.png was not
	// found, then default 16 color palette is used.
	//
	// Can be freely read and updated. Changes will be visible immediately.
	Palette [256]RGB

	DrawPalette    [256]byte
	DisplayPalette [256]byte

	ColorTransparency [256]bool
)

// RGB represents color
type RGB struct{ R, G, B byte }

func (r RGB) String() string {
	var rgb = int(r.R)<<16 + int(r.G)<<8 + int(r.B)
	return fmt.Sprintf("#%.6x", rgb) // avoid dependency on fmt inside the entire package
}
