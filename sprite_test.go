// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	_ "embed"
	"github.com/elgopher/pi/pitest"
	"testing"

	"github.com/elgopher/pi"
)

var (
	//go:embed internal/test/stretch/sprite.png
	spritePNG []byte
	//go:embed internal/test/stretch/sprite-3x3.png
	sprite3x3PNG []byte
	//go:embed internal/test/stretch/sprite-6x6.png
	sprite6x6PNG []byte
	//go:embed internal/test/stretch/sprite-6x3.png
	sprite6x3PNG []byte
	//go:embed internal/test/stretch/sprite-3x6.png
	sprite3x6PNG []byte
	//go:embed internal/test/stretch/sprite-4x3.png
	sprite4x3PNG []byte
	//go:embed internal/test/stretch/sprite-5x3.png
	sprite5x3PNG []byte
	//go:embed internal/test/stretch/sprite-3x4.png
	sprite3x4PNG []byte
	//go:embed internal/test/stretch/sprite-3x5.png
	sprite3x5PNG []byte
	//go:embed internal/test/stretch/sprite-1x1.png
	sprite1x1PNG []byte
	//go:embed internal/test/stretch/sprite-2x1.png
	sprite2x1PNG []byte
)

func TestStretch(t *testing.T) {
	t.Run("inside screen", func(t *testing.T) {
		tests := map[string]struct {
			dw, dh int
			png    []byte
		}{
			"3x3": {
				dw: 3, dh: 3,
				png: sprite3x3PNG,
			},
			"6x6": {
				dw: 6, dh: 6,
				png: sprite6x6PNG,
			},
			"6x3": {
				dw: 6, dh: 3,
				png: sprite6x3PNG,
			},
			"3x6": {
				dw: 3, dh: 6,
				png: sprite3x6PNG,
			},
			"4x3": {
				dw: 4, dh: 3,
				png: sprite4x3PNG,
			},
			"5x3": {
				dw: 5, dh: 3,
				png: sprite5x3PNG,
			},
			"3x4": {
				dw: 3, dh: 4,
				png: sprite3x4PNG,
			},
			"3x5": {
				dw: 3, dh: 5,
				png: sprite3x5PNG,
			},
			"1x1": {
				dw: 1, dh: 1,
				png: sprite1x1PNG,
			},
			"2x1": {
				dw: 2, dh: 1,
				png: sprite2x1PNG,
			},
		}
		for testName, testCase := range tests {
			t.Run(testName, func(t *testing.T) {
				pi.SetScreenSize(8, 8)
				pi.Cls()
				pi.Palette = pi.DecodePalette(spritePNG)
				sprite := pi.SpriteFrom(pi.DecodeCanvas(spritePNG), 1, 1, 3, 3)
				// when
				pi.Stretch(sprite, 1, 1, testCase.dw, testCase.dh)
				// then
				expected := pi.DecodeCanvas(testCase.png)
				pitest.AssertSurfaceEqual(t, expected, pi.Screen())
			})
		}
	})

	// temporary test
	dst := pi.NewCanvas(16, 16)
	pi.SetDrawTarget(dst)

	src := pi.NewCanvas(8, 8)
	src.Clear(7)

	spr := pi.CanvasSprite(src)

	pi.Stretch(spr, 0, 0, 8, 8)
	pi.Stretch(spr, -1, 0, 8, 8)
	pi.Stretch(spr, 0, -1, 8, 8)
	pi.Stretch(spr, 16, 0, 8, 8)
	pi.Stretch(spr, 0, 16, 8, 8)

	pi.Stretch(spr.WithFlipX(true), 0, 0, 8, 8)
	pi.Stretch(spr.WithFlipY(true), 0, 0, 8, 8)

	pi.Stretch(spr.WithSize(0, 0), 0, 0, 8, 8)
}
