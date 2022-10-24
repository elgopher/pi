// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/mem"
)

func TestSset(t *testing.T) {
	col := byte(2)

	t.Run("should set color of pixel in sprite sheet", func(t *testing.T) {
		pi.SpriteSheetWidth = 8
		pi.SpriteSheetHeight = 8
		pi.Resources = embed.FS{}
		pi.MustBoot()
		// when
		pi.Sset(2, 1, col)
		// then
		assert.Equal(t, col, mem.SpriteSheetData[10])
	})

	t.Run("should not set pixel outside the sprite sheet", func(t *testing.T) {
		pi.SpriteSheetWidth = 8
		pi.SpriteSheetHeight = 8
		pi.MustBoot()

		emptySheet := make([]byte, len(mem.SpriteSheetData))

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
				assert.Equal(t, emptySheet, mem.SpriteSheetData)
			})
		}
	})
}

func TestSget(t *testing.T) {
	t.Run("should get color of pixel", func(t *testing.T) {
		pi.SpriteSheetWidth = 8
		pi.SpriteSheetHeight = 8
		pi.Resources = embed.FS{}
		pi.MustBoot()
		col := byte(7)
		pi.Sset(1, 1, col)
		// expect
		assert.Equal(t, col, pi.Sget(1, 1))
	})

	t.Run("should get color 0 if outside the sprite sheet", func(t *testing.T) {
		pi.SpriteSheetWidth = 8
		pi.SpriteSheetHeight = 8
		pi.Resources = embed.FS{}
		pi.MustBoot()
		for i := 0; i < len(mem.SpriteSheetData); i++ {
			mem.SpriteSheetData[i] = 7
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
