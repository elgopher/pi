// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"strconv"
	"testing"

	"github.com/elgopher/pi"
	"github.com/stretchr/testify/assert"
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
}

func TestTime(t *testing.T) {
	t.Run("should return 0.0 when game was not run", func(t *testing.T) {
		assert.Zero(t, pi.Time())
	})
}
