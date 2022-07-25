// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"strconv"
	"testing"
	"testing/fstest"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoot(t *testing.T) {
	invalidSpriteSheetSizes := [...]int{
		0, 1, 7, 9,
	}

	t.Run("should return error if SpriteSheetWidth is not multiplication of 8", func(t *testing.T) {
		for _, width := range invalidSpriteSheetSizes {
			t.Run(strconv.Itoa(width), func(t *testing.T) {
				pi.Reset()
				pi.SpriteSheetWidth = width
				err := pi.Boot()
				assert.Error(t, err)
			})
		}
	})

	t.Run("should return error if SpriteSheetHeight is not multiplication of 8", func(t *testing.T) {
		for _, height := range invalidSpriteSheetSizes {
			t.Run(strconv.Itoa(height), func(t *testing.T) {
				pi.Reset()
				pi.SpriteSheetHeight = height
				err := pi.Boot()
				assert.Error(t, err)
			})
		}
	})

	invalidScreenSizes := [...]int{-2, -1, 0}

	t.Run("should return error when ScreenWidth is not greater than 0", func(t *testing.T) {
		for _, size := range invalidScreenSizes {
			t.Run(strconv.Itoa(size), func(t *testing.T) {
				pi.Reset()
				pi.ScreenWidth = size
				err := pi.Boot()
				assert.Error(t, err)
			})
		}
	})

	t.Run("should return error when ScreenHeight is not greater than 0", func(t *testing.T) {
		for _, size := range invalidScreenSizes {
			t.Run(strconv.Itoa(size), func(t *testing.T) {
				pi.Reset()
				pi.ScreenHeight = size
				err := pi.Boot()
				assert.Error(t, err)
			})
		}
	})

	t.Run("should use custom size sprite sheet when sprite-sheet.png was not found in resources", func(t *testing.T) {
		pi.Reset()
		pi.SpriteSheetWidth = 16
		pi.SpriteSheetHeight = 8
		allBlacks := [256]image.RGB{}
		pi.Palette = allBlacks
		// when
		err := pi.Boot()
		// then
		require.NoError(t, err)
		expectedSpriteSheetData := make([]byte, 16*8)
		assert.Equal(t, expectedSpriteSheetData, pi.SpriteSheetData)
		assert.Equal(t, allBlacks, pi.Palette)
	})

	t.Run("should load sprite-sheet.png", func(t *testing.T) {
		pi.Reset()
		pi.Resources = fstest.MapFS{
			"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet16x16},
		}
		// when
		err := pi.Boot()
		// then
		require.NoError(t, err)
		assert.Equal(t, 16, pi.SpriteSheetWidth)
		assert.Equal(t, 16, pi.SpriteSheetHeight)
		img := decodePNG(t, "internal/testimage/sprite-sheet-16x16.png")
		assert.Equal(t, img.Pixels, pi.SpriteSheetData)
		assert.Equal(t, img.Palette, pi.Palette)
	})

	t.Run("should reset draw state", func(t *testing.T) {
		pi.Reset()
		require.NoError(t, pi.Boot())
		pi.Color = 14
		pi.Camera(1, 2)
		pi.Clip(1, 2, 3, 4)
		// when
		err := pi.Boot()
		// then
		require.NoError(t, err)
		camX, camY := pi.CameraReset()
		assert.Zero(t, camX)
		assert.Zero(t, camY)
		x, y, w, h := pi.ClipReset()
		assert.Zero(t, x)
		assert.Zero(t, y)
		assert.Equal(t, pi.ScreenWidth, w)
		assert.Equal(t, pi.ScreenHeight, h)
		assert.Equal(t, byte(14), pi.Color)
	})
}

func TestTime(t *testing.T) {
	t.Run("should return 0.0 when game was not run", func(t *testing.T) {
		assert.Zero(t, pi.Time())
	})
}
