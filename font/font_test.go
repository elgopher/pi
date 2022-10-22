// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package font_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi/font"
)

var (
	//go:embed internal/invalid-width.png
	invalidWidthImage []byte

	//go:embed internal/invalid-height.png
	invalidHeightImage []byte

	//go:embed internal/valid.png
	validImage []byte

	//go:embed internal/empty.png
	emptyImage []byte
)

func TestLoadImageInto(t *testing.T) {
	t.Run("should return error when font png is invalid", func(t *testing.T) {
		var out [2048]byte
		err := font.Load(make([]byte, 0), out[:])
		require.Error(t, err)
	})

	t.Run("should return error for invalid image size", func(t *testing.T) {
		tests := map[string][]byte{
			"width not 128":  invalidWidthImage,
			"height not 128": invalidHeightImage,
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
		var out [1]byte
		err := font.Load(validImage, out[:])
		assert.Error(t, err)
	})

	t.Run("should override existing data", func(t *testing.T) {
		out := makeNotZeroSlice(2048, 1)
		// when
		err := font.Load(emptyImage, out)
		require.NoError(t, err)
		assert.Equal(t, make([]byte, 2048), out)
	})

	t.Run("should load pixels", func(t *testing.T) {
		out := make([]byte, 2048)
		// when
		err := font.Load(validImage, out)
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
