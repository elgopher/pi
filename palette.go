// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Palette swapping variables. Use them to generate new graphics by swapping colors.
var (
	// Pal (draw palette) contains mapping of colors used to replace color with
	// another one for all subsequent drawings. Affected functions are:
	// Pset, Spr, SprSize, SprSizeFlip, Circ, CircFill, Line, Rect and RectFill.
	//
	// The index of array is original color, the value is the color replacement.
	// For example,
	//
	//	pi.Pal[7] = 0
	//
	// will change color 7 to 0.
	Pal PalMapping

	// Pald (display palette), a mapping of colors used to replace color with
	// another one for the entire screen, at the end of a frame.
	//
	// The index of array is original color, the value is color replacement.
	// For example,
	//
	//   pi.Pald[7] = 0
	//
	// will change color 7 to 0.
	Pald PalMapping
)

type PalMapping [256]byte

// Reset resets all swapped colors to defaults (no swapping).
func (p *PalMapping) Reset() {
	*p = notSwappedPalette
}

// Palt contains information whether given color is transparent.
// If true then the color will not be drawn for next drawing operations.
// Color transparency is used by Spr, SprSize and SprSizeFlip.
//
// The index of array is a color number. For example,
//
//	pi.Palt[7] = true
//
// will make color 7 transparent.
var Palt Transparency = defaultTransparency

var notSwappedPalette [256]byte

func init() {
	for i := 0; i < 256; i++ {
		c := byte(i)
		notSwappedPalette[i] = c
		Pal[i] = c
		Pald[i] = c
	}
}

var defaultTransparency = Transparency{true}

type Transparency [256]bool

// Reset sets all transparent colors to false and makes color 0 transparent.
func (p *Transparency) Reset() {
	*p = defaultTransparency
}
