// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strings"

	"github.com/elgopher/pi/internal"
	"github.com/elgopher/pi/pimath"
)

// Palette is the global color palette used before the final rendering
// of the Screen contents to the game window. It maps Color values to RGB.
//
// The array index is the Color value, and the element is the RGB color,
// e.g., 0xFF00FF.
var Palette PaletteArray = defaultPalette()

// PaletteArray maps each Color to an RGB value.
//
// The array index is the Color, and the value is the RGB color,
// e.g., 0xFF00FF.
type PaletteArray [MaxColors]RGB

func (p PaletteArray) String() string {
	var s strings.Builder
	s.WriteString("{")
	for i := 0; i < len(p); i++ {
		s.WriteString(fmt.Sprintf("%d:%v, ", i, p[i]))
	}
	s.WriteString("}")
	return s.String()
}

// RGB represents a color with three components: R (Red), G (Green), and B (Blue).
//
// For example: black is 0x000000, white is 0xFFFFFF, green is 0x00FF00.
type RGB uint32

// FromRGB creates an RGB color from r, g, and b components.
func FromRGB(r, g, b uint8) RGB {
	return RGB(uint32(r)<<16 + uint32(g)<<8 + uint32(b))
}

// RGB returns the red, green, and blue components as 8-bit values.
func (rgb RGB) RGB() (r, g, b uint8) {
	r = uint8(rgb >> 16)
	g = uint8(rgb >> 8)
	b = uint8(rgb)
	return
}

// FromRGBf creates an RGB value from normalized floating-point components.
//
// The r, g, and b parameters should be in the range [0.0, 1.0].
// Values are clamped, converted to 8-bit, and combined into a single RGB value.
func FromRGBf(r, g, b float64) RGB {
	red := pimath.Clamp(r*255, 0, 255)
	green := pimath.Clamp(g*255, 0, 255)
	blue := pimath.Clamp(b*255, 0, 255)
	return RGB(uint32(red)<<16 + uint32(green)<<8 + uint32(blue))
}

// RGBf returns the red, green, and blue components as normalized float values.
//
// Each component is in the range [0.0, 1.0].
func (rgb RGB) RGBf() (r, g, b float64) {
	red, green, blue := rgb.RGB()
	r = float64(red) / 255
	g = float64(green) / 255
	b = float64(blue) / 255
	return
}

func (rgb RGB) String() string {
	return fmt.Sprintf("0x%06X", uint32(rgb))
}

// ResetPalette resets the palette to the default Picotron palette.
func ResetPalette() {
	Palette = defaultPalette()
}

func defaultPalette() PaletteArray {
	return PaletteArray{
		// Pico-8 colors:
		0x000000, // 0 - black
		0x1D2B53, // 1 - dark blue
		0x7E2553, // 2 - dark purple
		0x008751, // 3 - dark green
		0xAB5236, // 4 - brown
		0x5F574F, // 5 - dark gray
		0xC2C3C7, // 6 - light gray
		0xFFF1E8, // 7 - white
		0xFF004D, // 8 - red
		0xFFA300, // 9 - orange
		0xFFEC27, // 10 - yellow
		0x00E436, // 11 - green
		0x29ADFF, // 12 - blue
		0x83769C, // 13 - indigo
		0xFF77A8, // 14 - pink
		0xFFCCAA, // 15 - peach
		// Extra Picotron colors:
		0x2463B0, // 16 - true-blue
		0x00A5A1, // 17 - teal
		0x654688, // 18 - purple
		0x125359, // 19 - dark-teal
		0x742F29, // 20 - dark-brown
		0x452D32, // 21 - darker-grey
		0xA28879, // 22 - medium-grey
		0xFFACC5, // 23 - light-pink
		0xB9003E, // 24 - dark-red
		0xE26B02, // 25 - dark-orange
		0x95F042, // 26 - lime-green
		0x00B251, // 27 - medium-green
		0x64DFF6, // 28 - light-blue
		0xBD9ADF, // 29 - mauve
		0xE40DAB, // 30 - magenta
		0xFF8557, // 31 - peach
	}
}

var errToManyColors = fmt.Errorf(
	"PNG file has too many colors in indexed palette. "+
		"The maximum number is %d", MaxColors)

// DecodePalette extracts a palette from a PNG file.
//
//   - If the file uses an indexed palette, colors are read directly
//     with their assigned indices.
//   - Otherwise, all unique colors are read and assigned indices
//     automatically. Scanning starts from the top-left pixel
//     and proceeds line by line.
func DecodePalette(pngFile []byte) PaletteArray {
	p, err := DecodePaletteOrErr(pngFile)
	if err != nil {
		panic("DecodePalette failed: " + err.Error())
	}
	return p
}

// DecodePaletteOrErr works like DecodePalette but returns an error
// if the PNG file is invalid or contains too many colors.
func DecodePaletteOrErr(pngFile []byte) (PaletteArray, error) {
	stdImage, err := png.Decode(bytes.NewReader(pngFile))
	if err != nil {
		return PaletteArray{}, fmt.Errorf("PNG decoding failed: %w", err)
	}

	if indexedPalette, ok := stdImage.ColorModel().(color.Palette); ok {
		if len(indexedPalette) > MaxColors {
			return PaletteArray{}, errToManyColors
		}
		return convertIndexedPaletteToRGB(indexedPalette), nil
	}

	palette := internal.PaletteMaker[RGB]{}

	bounds := stdImage.Bounds()
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			c := stdImage.At(x, y)
			if palette.Add(c) {
				return PaletteArray{}, errToManyColors
			}
		}
	}

	return palette.Palette(), nil
}

func convertIndexedPaletteToRGB(indexedPalette color.Palette) PaletteArray {
	var p PaletteArray
	for i, col := range indexedPalette {
		r, g, b, _ := col.RGBA()
		r &= 0xff
		g &= 0xff
		b &= 0xff
		p[i] = RGB(r<<16 + g<<8 + b)
	}
	return p
}
