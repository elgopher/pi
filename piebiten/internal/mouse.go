// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/elgopher/pi/pimouse"
)

func (g *EbitenGame) updateMouse() {
	g.publishEventDown(ebiten.MouseButtonLeft, pimouse.Left)
	g.publishEventDown(ebiten.MouseButtonMiddle, "Middle")
	g.publishEventDown(ebiten.MouseButtonRight, pimouse.Right)
	g.publishEventDown(ebiten.MouseButton3, "3")
	g.publishEventDown(ebiten.MouseButton4, "4")

	g.publishEventUp(ebiten.MouseButtonLeft, pimouse.Left)
	g.publishEventUp(ebiten.MouseButtonMiddle, "Middle")
	g.publishEventUp(ebiten.MouseButtonRight, pimouse.Right)
	g.publishEventUp(ebiten.MouseButton3, "3")
	g.publishEventUp(ebiten.MouseButton4, "4")

	x, y := ebiten.CursorPosition()
	prev := g.mousePosition
	g.mousePosition.X = int((float64(x) - g.left) / g.scale)
	g.mousePosition.Y = int((float64(y) - g.top) / g.scale)

	mouseMovementDelta := g.mousePosition.Subtract(prev)

	if mouseMovementDelta.X != 0 || mouseMovementDelta.Y != 0 {
		event := pimouse.EventMove{
			Position: g.mousePosition,
			Previous: prev,
		}
		if !g.paused {
			pimouse.MoveTarget().Publish(event)
		}
		pimouse.MoveDebugTarget().Publish(event)
	}
}

func (g *EbitenGame) publishEventDown(button ebiten.MouseButton, key pimouse.Button) {
	if inpututil.IsMouseButtonJustPressed(button) {
		event := pimouse.EventButton{
			Type:   pimouse.EventButtonDown,
			Button: key,
		}
		if !g.paused {
			pimouse.ButtonTarget().Publish(event)
		}
		pimouse.ButtonDebugTarget().Publish(event)
	}
}

func (g *EbitenGame) publishEventUp(button ebiten.MouseButton, key pimouse.Button) {
	if inpututil.IsMouseButtonJustReleased(button) {
		event := pimouse.EventButton{
			Type:   pimouse.EventButtonUp,
			Button: key,
		}
		if !g.paused {
			pimouse.ButtonTarget().Publish(event)
		}
		pimouse.ButtonDebugTarget().Publish(event)
	}
}
