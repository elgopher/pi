// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package font_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/elgopher/pi/image"
	"github.com/elgopher/pi/internal/font"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed font_test.png
var testFont []byte

func TestLoadImageInto(t *testing.T) {
	t.Run("should return error for invalid image", func(t *testing.T) {
		tests := map[string]image.Image{
			"width not 128":            {Width: 1, Height: 128, Pixels: make([]byte, 128)},
			"height not 128":           {Width: 128, Height: 1, Pixels: make([]byte, 128)},
			"invalid number of pixels": {Width: 128, Height: 128, Pixels: make([]byte, 128)},
		}
		for name, img := range tests {
			t.Run(name, func(t *testing.T) {
				var out [2048]byte
				err := font.Load(img, out[:])
				assert.Error(t, err)
			})
		}
	})

	t.Run("should return error for invalid fontData", func(t *testing.T) {
		img := image.Image{Width: 128, Height: 128, Pixels: make([]byte, 128*128)}
		var out [1]byte
		err := font.Load(img, out[:])
		assert.Error(t, err)
	})

	t.Run("should override existing data", func(t *testing.T) {
		out := makeNotZeroSlice(2048, 1)
		emptyImage := image.Image{Width: 128, Height: 128, Pixels: make([]byte, 128*128)}
		// when
		err := font.Load(emptyImage, out)
		require.NoError(t, err)
		assert.Equal(t, make([]byte, 2048), out)
	})

	t.Run("should load pixels", func(t *testing.T) {
		out := make([]byte, 2048)
		img, err := image.DecodePNG(bytes.NewReader(testFont))
		require.NoError(t, err)
		// when
		err = font.Load(img, out)
		// then
		require.NoError(t, err)
		expectedChar0 := []byte{1, 2, 4, 8, 0x10, 0x20, 0x40, 0x80}
		assert.Equal(t, expectedChar0, out[0:8])
		expectedChar1 := []byte{0x80, 0x40, 0x20, 0x10, 8, 4, 2, 1}
		assert.Equal(t, expectedChar1, out[8:16])
		expectedChar2 := []byte{3, 7, 0xF, 0x1F, 0x3F, 0x7F, 0xFF, 0xFF}
		assert.Equal(t, expectedChar2, out[16:24])
		expectedChar15 := []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
		assert.Equal(t, expectedChar15, out[15*8:15*8+8])
		expectedChar16 := []byte{2, 4, 8, 0x10, 0x20, 0x40, 0x80, 1}
		assert.Equal(t, expectedChar16, out[16*8:16*8+8])
		expectedChar255 := []byte{0, 0, 0, 0, 0, 0, 0, 0x80}
		assert.Equal(t, expectedChar255, out[len(out)-8:])
	})
}

func makeNotZeroSlice(len int, fillWith byte) []byte {
	slice := make([]byte, len)
	for i := 0; i < len; i++ {
		slice[i] = fillWith
	}
	return slice
}
