// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
)

var allMouseButtons = []pi.MouseButton{pi.MouseLeft, pi.MouseMiddle, pi.MouseRight}

func TestMouseBtn(t *testing.T) {
	testMouseBtn(t, pi.MouseBtn)
}

func testMouseBtn(t *testing.T, mouseBtn func(button pi.MouseButton) bool) {
	t.Run("should return false for invalid button", func(t *testing.T) {
		assert.False(t, mouseBtn(-1))
		assert.False(t, mouseBtn(3))
	})

	t.Run("should not panic", func(t *testing.T) {
		for _, button := range allMouseButtons {
			mouseBtn(button)
		}
	})
}

func TestMouseBtnp(t *testing.T) {
	testMouseBtn(t, pi.MouseBtnp)
}

func TestMousePos(t *testing.T) {
	t.Run("should not panic", func(t *testing.T) {
		pi.MousePos()
	})
}
