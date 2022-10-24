// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package snap provides functions for taking screenshots.
package snap

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"

	"github.com/elgopher/pi/mem"
)

// Take takes a screenshot and saves it to temp dir.
//
// Take returns a filename. If something went wrong error is returned.
func Take() (string, error) {
	if runtime.GOOS == "js" {
		return "", fmt.Errorf("storing files does not work on js")
	}

	var palette color.Palette
	for _, col := range mem.DisplayPalette {
		rgb := mem.Palette[col]
		rgba := &color.NRGBA{R: rgb.R, G: rgb.G, B: rgb.B, A: 255}
		palette = append(palette, rgba)
	}

	size := image.Rectangle{Max: image.Point{X: mem.ScreenWidth, Y: mem.ScreenHeight}}
	img := image.NewPaletted(size, palette)

	copy(img.Pix, mem.ScreenData)

	file, err := os.CreateTemp("", "screenshot-*.png")
	if err != nil {
		return "", fmt.Errorf("error creating temp file for screenshot: %w", err)
	}

	if err = png.Encode(file, img); err != nil {
		return "", fmt.Errorf("error encoding screenshot into PNG file: %w", err)
	}

	if err = file.Close(); err != nil {
		return "", fmt.Errorf("error closing file: %w", err)
	}

	return file.Name(), nil
}
