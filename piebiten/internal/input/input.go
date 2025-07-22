// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package input

import (
	"github.com/elgopher/pi"
	"github.com/hajimehoshi/ebiten/v2"
)

type Backend struct {
	Paused     *bool
	LeftOffset *float64
	TopOffset  *float64
	Scale      *float64

	keys          []ebiten.Key
	mousePosition pi.Position
	gamepads      ebitenGamepads
}

func (g *Backend) Update() {
	g.updateMouse()
	g.updateKeyboard()
	g.gamepads.update()
}
