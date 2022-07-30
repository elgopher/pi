// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/elgopher/pi"
	"github.com/stretchr/testify/assert"
)

var allButtons = []pi.Button{pi.Left, pi.Right, pi.Up, pi.Down, pi.O, pi.X}

func TestBtn(t *testing.T) {
	testBtn(t, pi.Btn)
}

func testBtn(t *testing.T, btn func(pi.Button) bool) {
	t.Run("should return false for invalid button", func(t *testing.T) {
		assert.False(t, btn(-1))
		assert.False(t, btn(6))
	})

	t.Run("should not panic", func(t *testing.T) {
		for _, button := range allButtons {
			btn(button)
		}
	})
}

func TestBtnPlayer(t *testing.T) {
	testBtnPlayer(t, pi.BtnPlayer)
}

func testBtnPlayer(t *testing.T, btnPlayer func(pi.Button, int) bool) {
	testBtn(t, func(b pi.Button) bool {
		return btnPlayer(b, 0)
	})

	t.Run("should return false for invalid player", func(t *testing.T) {
		assert.False(t, btnPlayer(pi.X, -1))
		assert.False(t, btnPlayer(pi.X, 8))
	})

	t.Run("should not panic", func(t *testing.T) {
		for player := 0; player < 8; player++ {
			for _, button := range allButtons {
				btnPlayer(button, player)
			}
		}
	})
}

func TestBtnp(t *testing.T) {
	testBtn(t, pi.Btnp)
}

func TestBtnpPlayer(t *testing.T) {
	testBtnPlayer(t, pi.BtnpPlayer)
}

func TestBtnBits(t *testing.T) {
	t.Run("should not panic", func(t *testing.T) {
		pi.BtnBits()
	})
}

func TestBtnpBits(t *testing.T) {
	t.Run("should not panic", func(t *testing.T) {
		pi.BtnpBits()
	})
}
