// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package image

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

// DecodePNG decodes PNG with indexed color mode. Such image must have at most 256 colors.
func DecodePNG(reader io.Reader) (Image, error) {
	var img Image
	if reader == nil {
		return img, errors.New("nil reader")
	}

	stdImage, err := png.Decode(reader)
	if err != nil {
		return img, fmt.Errorf("error decoding PNG file: %w", err)
	}

	indexedPalette, ok := stdImage.ColorModel().(color.Palette)
	if !ok {
		return img, errors.New("image does not have indexed color mode")
	}

	bounds := stdImage.Bounds()

	return Image{
		Width:   bounds.Dx(),
		Height:  bounds.Dy(),
		Palette: convertPaletteToRGB(indexedPalette),
		Pix:     convertImageToPixels(stdImage, indexedPalette),
	}, nil
}

func convertPaletteToRGB(palette color.Palette) [256]RGB {
	var p [256]RGB
	for i, col := range palette {
		r, g, b, _ := col.RGBA()
		p[i] = RGB{R: uint8(r), G: uint8(g), B: uint8(b)}
	}
	return p
}

func convertImageToPixels(img image.Image, palette color.Palette) []byte {
	bounds := img.Bounds()
	pixels := make([]byte, bounds.Dx()*bounds.Dy())

	offset := 0
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			col := img.At(x, y)
			pixels[offset] = indexOfColorInPalette(col, palette)
			offset++
		}
	}

	return pixels
}

func indexOfColorInPalette(col color.Color, palette color.Palette) byte {
	for i, paletteCol := range palette {
		if paletteCol == col {
			return byte(i)
		}
	}

	return 0
}
