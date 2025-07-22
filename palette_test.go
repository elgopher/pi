// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	_ "embed"
	"github.com/elgopher/pi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	//go:embed internal/test/palette/indexed.png
	indexedPalettePNG []byte
	//go:embed internal/test/palette/rgb.png
	rgbPalettePNG []byte

	//go:embed internal/test/palette/indexed-65.png
	indexed65PalettePNG []byte
	//go:embed internal/test/palette/rgb-65.png
	rgb65PalettePNG []byte

	//go:embed internal/test/palette/indexed-64.png
	indexed64PalettePNG []byte
	//go:embed internal/test/palette/rgb-64.png
	rgb64PalettePNG []byte
)

func TestDecodePaletteOrErr(t *testing.T) {
	t.Run("should return error when palette has too many colors", func(t *testing.T) {
		tests := map[string][]byte{
			"indexed": indexed65PalettePNG,
			"rgb":     rgb65PalettePNG,
		}
		for testName, file := range tests {
			t.Run(testName, func(t *testing.T) {
				palette, err := pi.DecodePaletteOrErr(file)
				require.ErrorContains(t, err, "too many colors")
				assert.Zero(t, palette)
			})
		}
	})

	t.Run("should not return error when palette has maximum number of colors", func(t *testing.T) {
		tests := map[string][]byte{
			"indexed": indexed64PalettePNG,
			"rgb":     rgb64PalettePNG,
		}
		for testName, file := range tests {
			t.Run(testName, func(t *testing.T) {
				_, err := pi.DecodePaletteOrErr(file)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("should decode palette", func(t *testing.T) {
		expected := pi.PaletteArray{0x000000, 0x1900ff, 0xff00b0, 0xff2600}
		tests := map[string][]byte{
			"indexed": indexedPalettePNG,
			"rgb":     rgbPalettePNG,
		}
		for testName, file := range tests {
			t.Run(testName, func(t *testing.T) {
				palette, err := pi.DecodePaletteOrErr(file)
				require.NoError(t, err)
				assert.Equal(t, expected, palette)
			})
		}
	})
}
