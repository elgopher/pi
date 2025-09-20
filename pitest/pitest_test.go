// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pitest_test

import (
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pitest"
	"github.com/stretchr/testify/assert"
)

func TestAssertSurfaceEqual(t *testing.T) {
	t.Run("should not fail when surfaces are equal", func(t *testing.T) {
		surface1 := pi.NewSurface[int](2, 3)
		surface1.SetAll(0, 1, 2, 3, 4, 5)
		surface2 := surface1.Clone()
		mockT := new(testing.T)
		// when
		actual := pitest.AssertSurfaceEqual(mockT, surface1, surface2)
		// then
		assert.False(t, mockT.Failed())
		assert.True(t, actual)
	})

	t.Run("should fail when surfaces are not equal", func(t *testing.T) {
		tests := map[string]struct {
			surface1, surface2 pi.Surface[int]
		}{
			"different data": {
				surface1: func() pi.Surface[int] {
					surface := pi.NewSurface[int](1, 1)
					surface.SetAll(1, 2)
					return surface
				}(),
				surface2: pi.NewSurface[int](1, 1),
			},
			"different width": {
				surface1: pi.NewSurface[int](1, 1),
				surface2: pi.NewSurface[int](2, 1),
			},
			"different height": {
				surface1: pi.NewSurface[int](1, 1),
				surface2: pi.NewSurface[int](1, 2),
			},
		}
		for testName, testCase := range tests {
			t.Run(testName, func(t *testing.T) {
				mockT := new(testing.T)
				// when
				actual := pitest.AssertSurfaceEqual(mockT, testCase.surface1, testCase.surface2)
				// then
				assert.True(t, mockT.Failed())
				assert.False(t, actual)
			})
		}
	})
}

func TestAssertSpriteEqual(t *testing.T) {
	t.Run("should not fail when sprites are equal", func(t *testing.T) {
		canvas := pi.NewCanvas(1, 1)
		canvas.SetAll(1, 2)
		sprite1 := pi.CanvasSprite(canvas)
		sprite2 := pi.CanvasSprite(canvas)
		mockT := new(testing.T)
		// when
		pitest.AssertSpriteEqual(mockT, sprite1, sprite2)
		// then
		assert.False(t, mockT.Failed())
	})

	t.Run("should fail when sprites are not equal", func(t *testing.T) {
		var data = []pi.Color{1, 2}
		canvas := pi.NewCanvas(1, 1)
		canvas.SetData(data)

		tests := map[string]struct {
			sprite1, sprite2 pi.Sprite
		}{
			"different canvas": {
				sprite1: func() pi.Sprite {
					return pi.CanvasSprite(canvas)
				}(),
				sprite2: func() pi.Sprite {
					canvas := pi.NewCanvas(1, 1)
					canvas.SetData(data)
					return pi.CanvasSprite(canvas)
				}(),
			},
			"different area": {
				sprite1: func() pi.Sprite {
					return pi.CanvasSprite(canvas)
				}(),
				sprite2: func() pi.Sprite {
					return pi.SpriteFrom(canvas, 1, 1, 0, 0)
				}(),
			},
			"different flipX": {
				sprite1: func() pi.Sprite {
					return pi.CanvasSprite(canvas).WithFlipX(true)
				}(),
				sprite2: func() pi.Sprite {
					return pi.CanvasSprite(canvas)
				}(),
			},
			"different flipY": {
				sprite1: func() pi.Sprite {
					return pi.CanvasSprite(canvas).WithFlipY(true)
				}(),
				sprite2: func() pi.Sprite {
					return pi.CanvasSprite(canvas)
				}(),
			},
		}
		for testName, testCase := range tests {
			t.Run(testName, func(t *testing.T) {
				mockT := new(testing.T)
				// when
				pitest.AssertSpriteEqual(mockT, testCase.sprite1, testCase.sprite2)
				// then
				assert.True(t, mockT.Failed())
			})
		}
	})
}
