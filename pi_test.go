package pi_test

import (
	"strconv"
	"testing"

	"github.com/elgopher/pi"
	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	t.Run("should fail if SpriteSheetWidth is not multiplication of 8", func(t *testing.T) {
		tests := [...]int{
			0, 1, 7, 9,
		}
		for _, width := range tests {
			t.Run(strconv.Itoa(width), func(t *testing.T) {
				pi.SpriteSheetWidth = width
				err := pi.Boot()
				assert.Error(t, err)
			})
		}
	})
	t.Run("should fail if SpriteSheetHeight is not multiplication of 8", func(t *testing.T) {
		tests := [...]int{
			0, 1, 7, 9,
		}
		for _, height := range tests {
			t.Run(strconv.Itoa(height), func(t *testing.T) {
				pi.SpriteSheetHeight = height
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
