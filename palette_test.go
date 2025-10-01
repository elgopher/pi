// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi"
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

func TestFromRGB(t *testing.T) {
	// when
	rgb := pi.FromRGB(1, 128, 255)
	// then
	r, g, b := rgb.RGB()
	assert.Equal(t, uint8(1), r)
	assert.Equal(t, uint8(128), g)
	assert.Equal(t, uint8(255), b)
}

func TestFromRGBf(t *testing.T) {
	// when
	rgb := pi.FromRGBf(0.004, 0.5, 1.0)
	// then
	rf, gf, bf := rgb.RGBf()
	assert.InDelta(t, 0.004, rf, 0.01)
	assert.InDelta(t, 0.5, gf, 0.01)
	assert.InDelta(t, 1.0, bf, 0.01)
}

func TestRGB_String(t *testing.T) {
	rgb := pi.RGB(0x304050)
	assert.Equal(t, "0x304050", rgb.String())
}

func TestPaletteArray_String(t *testing.T) {
	var p pi.PaletteArray
	p[0] = pi.RGB(0x010101)
	p[1] = pi.RGB(0x020202)
	p[63] = pi.RGB(0x636363)
	// when
	actual := p.String()
	// then
	assert.Equal(t,
		"{0:0x010101, 1:0x020202, 2:0x000000, 3:0x000000, 4:0x000000, 5:0x000000, 6:0x000000, 7:0x000000, 8:0x000000, "+
			"9:0x000000, 10:0x000000, 11:0x000000, 12:0x000000, 13:0x000000, 14:0x000000, 15:0x000000, 16:0x000000, "+
			"17:0x000000, 18:0x000000, 19:0x000000, 20:0x000000, 21:0x000000, 22:0x000000, 23:0x000000, 24:0x000000, "+
			"25:0x000000, 26:0x000000, 27:0x000000, 28:0x000000, 29:0x000000, 30:0x000000, 31:0x000000, 32:0x000000, "+
			"33:0x000000, 34:0x000000, 35:0x000000, 36:0x000000, 37:0x000000, 38:0x000000, 39:0x000000, 40:0x000000, "+
			"41:0x000000, 42:0x000000, 43:0x000000, 44:0x000000, 45:0x000000, 46:0x000000, 47:0x000000, 48:0x000000, "+
			"49:0x000000, 50:0x000000, 51:0x000000, 52:0x000000, 53:0x000000, 54:0x000000, 55:0x000000, 56:0x000000, "+
			"57:0x000000, 58:0x000000, 59:0x000000, 60:0x000000, 61:0x000000, 62:0x000000, 63:0x636363, }",
		actual)
}
