// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package mem

var (
	SystemFont = Font{
		Width:        4,
		WidthSpecial: 8,
		Height:       6,
	}

	CustomFont = Font{
		Width:        4,
		WidthSpecial: 8,
		Height:       6,
	}
)

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
