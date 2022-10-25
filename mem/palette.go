// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package mem

var (
	// Palette has all colors available in the game. Up to 256.
	// Palette is taken from loaded sprite sheet (which must be
	// a PNG file with indexed color mode). If sprite-sheet.png was not
	// found, then default 16 color palette is used.
	//
	// Can be freely read and updated. Changes will be visible immediately.
	Palette [256]RGB

	// DrawPalette contains mapping of colors used to replace color with
	// another one for all subsequent drawings.
	//
	// The index of array is original color, the value is color replacement.
	DrawPalette [256]byte

	// DisplayPalette contains mapping of colors used to replace color with
	// another one for the entire screen, at the end of a frame
	//
	// The index of array is original color, the value is color replacement.
	DisplayPalette [256]byte

	// ColorTransparency contains information whether given color is transparent.
	//
	// The index of array is a color number.
	ColorTransparency [256]bool
)

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