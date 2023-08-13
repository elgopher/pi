// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/image"
)

func TestRGB_String(t *testing.T) {
	tests := map[string]image.RGB{
		"#000000": {},
		"#FFFFFF": {0xFF, 0xFF, 0xFF},
		"#012345": {0x01, 0x23, 0x45},
		"#6789AB": {0x67, 0x89, 0xAB},
		"#CDEF01": {0xCD, 0xEF, 0x01},
	}

	for expected, rgb := range tests {
		assert.Equal(t, expected, rgb.String())
	}
}

func TestImage_String(t *testing.T) {
	t.Run("should convert small image to string", func(t *testing.T) {
		img := image.Image{
			Width: 2, Height: 1,
			Palette: [256]image.RGB{{1, 1, 1}, {2, 2, 2}},
			Pix:     make([]byte, 2),
		}

		actual := img.String()
		expected := "{width:2, height:1, palette: (256)[#010101 #020202 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 #000000 ...], pix:[0 0]}"
		assert.Equal(t, expected, actual)
	})

	t.Run("should convert big image to string", func(t *testing.T) {
		img := image.Image{
			Pix: make([]byte, 100*100), // 10K bytes
		}
		actual := img.String()
		assert.True(t, len(actual) < 2500)
	})
}
