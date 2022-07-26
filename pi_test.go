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

	t.Run("should initialize screen data", func(t *testing.T) {
		pi.Reset()
		pi.ScreenWidth = 2
		pi.ScreenHeight = 3
		// when
		err := pi.Boot()
		// then
		require.NoError(t, err)
		assert.Equal(t, make([]byte, 6), pi.ScreenData)
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

	t.Run("changing the user parameters after Boot should not ends up in a panic", func(t *testing.T) {
		pi.Reset()
		pi.ScreenWidth = 8
		pi.ScreenHeight = 8
		pi.SpriteSheetWidth = 8
		pi.SpriteSheetHeight = 8
		pi.BootOrPanic()
		// when
		pi.ScreenWidth = 1
		pi.ScreenHeight = 1
		pi.SpriteSheetWidth = 1
		pi.SpriteSheetHeight = 1
		// then
		assert.NotPanics(t, func() {
			pi.Pset(1, 1)
			pi.Pget(1, 1)
			pi.Sset(1, 1, 7)
			pi.Sget(1, 1)
			pi.Spr(1, 1, 1)
			pi.SprSize(1, 1, 1, 1.0, 1.0)
			pi.SprSizeFlip(1, 1, 1, 1.0, 1.0, true, true)
			pi.Cls()
			pi.ClsCol(7)
		})
	})

	t.Run("should use sprite-sheet size loaded from sprite-sheet.png", func(t *testing.T) {
		pi.SpriteSheetWidth = 32 // 2x times bigger than actual sprite-sheet.png width
		pi.SpriteSheetHeight = 32
		pi.Resources = fstest.MapFS{
			"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet16x16},
		}
		pi.BootOrPanic()
		assert.NotPanics(t, func() {
			pi.Spr(4, 0, 0) // sprite-sheet.png has only 4 sprites (from 0 to 3)
			pi.SprSize(4, 0, 0, 1.0, 1.0)
			pi.SprSizeFlip(4, 0, 0, 1.0, 1.0, false, false)
			pi.Pset(16, 16) // sprite-sheet.png is only 16x16 pixels (0..15)
			pi.Pget(16, 16)
		})
	})
}

func TestTime(t *testing.T) {
	t.Run("should return 0.0 when game was not run", func(t *testing.T) {
		assert.Zero(t, pi.Time())
	})
}
