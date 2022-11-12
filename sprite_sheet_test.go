// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
)

func TestUseEmptySpriteSheet(t *testing.T) {
	invalidSpriteSheetSizes := [...]int{
		0, 1, 7, 9,
	}

	t.Run("should panic if SpriteSheetWidth is not multiple of 8", func(t *testing.T) {
		for _, width := range invalidSpriteSheetSizes {
			t.Run(strconv.Itoa(width), func(t *testing.T) {
				pi.Reset()
				assert.Panics(t, func() {
					pi.UseEmptySpriteSheet(width, pi.SprSheet().Height())
				})
			})
		}
	})

	t.Run("should panic if SpriteSheetHeight is not multiple of 8", func(t *testing.T) {
		for _, height := range invalidSpriteSheetSizes {
			t.Run(strconv.Itoa(height), func(t *testing.T) {
				pi.Reset()
				assert.Panics(t, func() {
					pi.UseEmptySpriteSheet(pi.SprSheet().Width(), height)
				})
			})
		}
	})

	t.Run("should use custom size sprite sheet", func(t *testing.T) {
		pi.Reset()
		// when
		pi.UseEmptySpriteSheet(16, 8)
		// then
		expectedSpriteSheetData := make([]byte, 16*8)
		assert.Equal(t, expectedSpriteSheetData, pi.SprSheet().Pix())
	})
}

func TestSset(t *testing.T) {
	col := byte(2)

	t.Run("should set color of pixel in sprite sheet", func(t *testing.T) {
		pi.UseEmptySpriteSheet(8, 8)
		// when
		pi.Sset(2, 1, col)
		// then
		assert.Equal(t, col, pi.SprSheet().Pix()[10])
	})

	t.Run("should not set pixel outside the sprite sheet", func(t *testing.T) {
		pi.UseEmptySpriteSheet(8, 8)

		emptySheet := make([]byte, len(pi.SprSheet().Pix()))

		tests := []struct{ X, Y int }{
			{-1, 0},
			{0, -1},
			{8, 0},
			{0, 8},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				// when
				pi.Sset(coords.X, coords.Y, col)
				// then
				assert.Equal(t, emptySheet, pi.SprSheet().Pix())
			})
		}
	})
}

func TestSget(t *testing.T) {
	t.Run("should get color of pixel", func(t *testing.T) {
		pi.UseEmptySpriteSheet(8, 8)
		col := byte(7)
		pi.Sset(1, 1, col)
		// expect
		assert.Equal(t, col, pi.Sget(1, 1))
	})

	t.Run("should get color 0 if outside the sprite sheet", func(t *testing.T) {
		pi.UseEmptySpriteSheet(8, 8)
		pixels := pi.SprSheet().Pix()
		for i := 0; i < len(pixels); i++ {
			pixels[i] = 7
		}

		tests := []struct{ X, Y int }{
			{-1, 0},
			{0, -1},
			{8, 0},
			{0, 8},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				// when
				actual := pi.Sget(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})
}
