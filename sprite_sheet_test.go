// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
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
