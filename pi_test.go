// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	_ "embed"
	"strconv"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/image"
)

//go:embed internal/testimage/custom-font.png
var customFont []byte

func TestReset(t *testing.T) {
	t.Run("should reset sprite sheet", func(t *testing.T) {
		pi.UseEmptySpriteSheet(8, 8)
		// when
		pi.Reset()
		// then
		sprSheet := pi.SprSheet()
		assert.Equal(t, 128, sprSheet.W)
		assert.Equal(t, 128, sprSheet.H)
		assert.Equal(t, make([]byte, 16384), sprSheet.Pix)
	})

	t.Run("should reset screen", func(t *testing.T) {
		pi.SetScreenSize(256, 256)
		// when
		pi.Reset()
		// then
		scr := pi.Scr()
		assert.Equal(t, 128, scr.W)
		assert.Equal(t, 128, scr.H)
		assert.Equal(t, make([]byte, 16384), scr.Pix)
		assert.Zero(t, scr.Camera)
		assert.Equal(t, pi.Region{W: 128, H: 128}, scr.Clip)
	})

	t.Run("should reset palette", func(t *testing.T) {
		color := image.RGB{R: 0xff, G: 0xff, B: 0xff}
		pi.Palette[0] = color
		pi.Reset()
		assert.NotEqual(t, color, pi.Palette[0])
	})

	t.Run("should reset display palette", func(t *testing.T) {
		pi.DisplayPalette[0] = 255
		pi.Reset()
		assert.NotEqual(t, 255, pi.DisplayPalette[0])
	})

	t.Run("should reset draw palette", func(t *testing.T) {
		pi.DrawPalette[0] = 255
		pi.Reset()
		assert.NotEqual(t, 255, pi.DrawPalette[0])
	})

	t.Run("should reset palette transparency", func(t *testing.T) {
		pi.ColorTransparency[0] = false
		pi.Reset()
		assert.True(t, pi.ColorTransparency[0])
	})

	t.Run("should reset system font", func(t *testing.T) {
		before := pi.SystemFont().Data[0]
		after := before + 1
		pi.SystemFont().Data[0] = after
		// when
		pi.Reset()
		// then
		assert.Equal(t, before, pi.SystemFont().Data[0])
	})

	t.Run("should reset custom font", func(t *testing.T) {
		pi.CustomFont().Data[0] = 1
		pi.SetCustomFontWidth(0)
		pi.SetCustomFontSpecialWidth(0)
		pi.SetCustomFontHeight(0)
		// when
		pi.Reset()
		// then
		expected := pi.Font{
			Data:         make([]byte, 8*256),
			Width:        4,
			SpecialWidth: 8,
			Height:       6,
		}
		assert.Equal(t, expected, pi.CustomFont())
	})

	t.Run("should reset callbacks", func(t *testing.T) {
		pi.Draw = func() {}
		pi.Update = func() {}
		pi.Reset()
		assert.Nil(t, pi.Draw)
		assert.Nil(t, pi.Update)
	})
}

func TestUseEmptySpriteSheet(t *testing.T) {
	invalidSpriteSheetSizes := [...]int{
		0, 1, 7, 9,
	}

	t.Run("should panic if SpriteSheetWidth is not multiplication of 8", func(t *testing.T) {
		for _, width := range invalidSpriteSheetSizes {
			t.Run(strconv.Itoa(width), func(t *testing.T) {
				pi.Reset()
				assert.Panics(t, func() {
					pi.UseEmptySpriteSheet(width, pi.SprSheet().H)
				})
			})
		}
	})

	t.Run("should panic if SpriteSheetHeight is not multiplication of 8", func(t *testing.T) {
		for _, height := range invalidSpriteSheetSizes {
			t.Run(strconv.Itoa(height), func(t *testing.T) {
				pi.Reset()
				assert.Panics(t, func() {
					pi.UseEmptySpriteSheet(pi.SprSheet().W, height)
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
		assert.Equal(t, expectedSpriteSheetData, pi.SprSheet().Pix)
	})
}

func TestSetScreenSize(t *testing.T) {
	invalidScreenSizes := [...]int{-2, -1, 0}

	t.Run("should panic when ScreenWidth is not greater than 0", func(t *testing.T) {
		for _, size := range invalidScreenSizes {
			t.Run(strconv.Itoa(size), func(t *testing.T) {
				pi.Reset()
				assert.Panics(t, func() {
					pi.SetScreenSize(size, pi.Scr().H)
				})
			})
		}
	})

	t.Run("should panic when ScreenHeight is not greater than 0", func(t *testing.T) {
		for _, size := range invalidScreenSizes {
			t.Run(strconv.Itoa(size), func(t *testing.T) {
				pi.Reset()
				assert.Panics(t, func() {
					pi.SetScreenSize(pi.Scr().W, size)
				})
			})
		}
	})

	t.Run("should initialize screen data", func(t *testing.T) {
		pi.Reset()
		// when
		pi.SetScreenSize(2, 3)
		// then
		assert.Equal(t, make([]byte, 6), pi.Scr().Pix)
	})
}

func TestLoad(t *testing.T) {
	const color = 7

	t.Run("should load sprite-sheet.png", func(t *testing.T) {
		pi.Reset()
		// when
		pi.Load(fstest.MapFS{
			"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet16x16},
		})
		// then
		assert.Equal(t, 16, pi.SprSheet().W)
		assert.Equal(t, 16, pi.SprSheet().H)
		img := decodePNG(t, "internal/testimage/sprite-sheet-16x16.png")
		assert.Equal(t, img.Pixels, pi.SprSheet().Pix)
		assert.Equal(t, img.Palette, pi.Palette)
	})

	t.Run("should load custom-font.png", func(t *testing.T) {
		pi.Reset()
		// when
		pi.Load(fstest.MapFS{
			"custom-font.png": &fstest.MapFile{Data: customFont},
		})
		// then
		assert.Equal(t, uint8(0xf), pi.CustomFont().Data[0])
	})

	t.Run("should use sprite-sheet size loaded from sprite-sheet.png", func(t *testing.T) {
		pi.UseEmptySpriteSheet(32, 32) // 2x times bigger than actual sprite-sheet.png width
		pi.Load(fstest.MapFS{
			"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet16x16},
		})
		assert.NotPanics(t, func() {
			pi.Spr(4, 0, 0) // sprite-sheet.png has only 4 sprites (from 0 to 3)
			pi.SprSize(4, 0, 0, 1.0, 1.0)
			pi.SprSizeFlip(4, 0, 0, 1.0, 1.0, false, false)
			pi.Pset(16, 16, color) // sprite-sheet.png is only 16x16 pixels (0..15)
			pi.Pget(16, 16)
		})
	})
}

func TestTime(t *testing.T) {
	t.Run("should return 0.0 when game was not run", func(t *testing.T) {
		assert.Zero(t, pi.Time())
	})
}
