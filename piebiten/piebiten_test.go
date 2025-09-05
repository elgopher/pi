// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package piebiten_test

import (
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/piebiten/internal/ebitentesting"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ebitentesting.MainWithRunLoop(m)
}

func TestCopyCanvasToEbitenImage(t *testing.T) {
	t.Run("should copy canvas to Ebiten image", func(t *testing.T) {
		canvas := pi.NewCanvas(3, 2)
		pi.Palette[0] = 0xFF0000
		pi.Palette[1] = 0x00FF00
		pi.Palette[2] = 0x0000FF
		pi.Palette[3] = 0x00FFFF
		pi.Palette[4] = 0xFFFF00
		pi.Palette[5] = 0xFF00FF
		canvas.SetAll(
			0, 1, 2,
			3, 4, 5,
		)
		img := ebiten.NewImage(3, 2)
		// when
		piebiten.CopyCanvasToEbitenImage(canvas, img)
		// then
		bounds := img.Bounds()
		out := make([]byte, bounds.Dx()*bounds.Dy()*4)
		img.ReadPixels(out)
		for i := 0; i < 6; i++ {
			pix := out[i*4 : i*4+4]
			rgb := pi.FromRGB(pix[0], pix[1], pix[2])
			assert.Equal(t, pi.Palette[i], rgb, "pixel %d", i)
			assert.Equal(t, uint8(0xff), pix[3], "alpha for pixel %d", i)
		}
	})

	t.Run("should take into account palette mapping", func(t *testing.T) {
		canvas := pi.NewCanvas(1, 1)
		pi.Palette[0] = 0x000000
		pi.Palette[1] = 0xFFFFFF
		pi.PaletteMapping[0] = 1 // map 0 to 1 (0x000000 to 0xFFFFFF)
		img := ebiten.NewImage(1, 1)
		// when
		piebiten.CopyCanvasToEbitenImage(canvas, img)
		// then
		bounds := img.Bounds()
		out := make([]byte, bounds.Dx()*bounds.Dy()*4)
		img.ReadPixels(out)
		pix := out[0:4]
		rgb := pi.FromRGB(pix[0], pix[1], pix[2])
		assert.Equal(t, pi.RGB(0xFFFFFF), rgb)
	})
}
