// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/mem"
)

func updateController() {
	for player := 0; player < 8; player++ {
		getController(player).update(player)
	}
}

func getController(player int) *controller {
	c := controller{&mem.Controllers[player]}
	return &c
}

type controller struct {
	*mem.Controller
}

func (c *controller) update(player int) {
	c.updateDirections(player)
	c.updateFireButtons(player)
}

func (c *controller) updateDirections(player int) {
	gamepadID := ebiten.GamepadID(player)

	axisX := ebiten.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal)
	axisY := ebiten.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxisLeftStickVertical)

	if axisX < -0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftLeft) ||
		isKeyboardPressed(player, mem.ControllerLeft) {
		c.BtnDuration[mem.ControllerLeft] += 1
		c.BtnDuration[mem.ControllerRight] = 0
	} else if axisX > 0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftRight) ||
		isKeyboardPressed(player, mem.ControllerRight) {
		c.BtnDuration[mem.ControllerRight] += 1
		c.BtnDuration[mem.ControllerLeft] = 0
	} else {
		c.BtnDuration[mem.ControllerRight] = 0
		c.BtnDuration[mem.ControllerLeft] = 0
	}

	if axisY < -0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftTop) ||
		isKeyboardPressed(player, mem.ControllerUp) {
		c.BtnDuration[mem.ControllerUp] += 1
		c.BtnDuration[mem.ControllerDown] = 0
	} else if axisY > 0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftBottom) ||
		isKeyboardPressed(player, mem.ControllerDown) {
		c.BtnDuration[mem.ControllerDown] += 1
		c.BtnDuration[mem.ControllerUp] = 0
	} else {
		c.BtnDuration[mem.ControllerUp] = 0
		c.BtnDuration[mem.ControllerDown] = 0
	}
}

func (c *controller) updateFireButtons(player int) {
	gamepadID := ebiten.GamepadID(player)

	if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightBottom) ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightTop) ||
		isKeyboardPressed(player, mem.ControllerO) {
		c.BtnDuration[mem.ControllerO] += 1
	} else {
		c.BtnDuration[mem.ControllerO] = 0
	}

	if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightRight) ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightLeft) ||
		isKeyboardPressed(player, mem.ControllerX) {
		c.BtnDuration[mem.ControllerX] += 1
	} else {
		c.BtnDuration[mem.ControllerX] = 0
	}
}

// first array is player, then π key, then slice of Ebitengine keys.
var keyboardMapping = [...][6][]ebiten.Key{
	// player0:
	{
		{ebiten.KeyLeft},                        // left
		{ebiten.KeyRight},                       // right
		{ebiten.KeyUp},                          // up
		{ebiten.KeyDown},                        // down
		{ebiten.KeyZ, ebiten.KeyC, ebiten.KeyN}, // o
		{ebiten.KeyX, ebiten.KeyV, ebiten.KeyM}, // x
	},
	// player1:
	{
		{ebiten.KeyS},         // left
		{ebiten.KeyF},         // right
		{ebiten.KeyE},         // up
		{ebiten.KeyD},         // down
		{ebiten.KeyShiftLeft}, // o
		{ebiten.KeyTab, ebiten.KeyW, ebiten.KeyQ, ebiten.KeyA}, // x
	},
}

func isKeyboardPressed(player int, button int) bool {
	if player >= len(keyboardMapping) {
		return false
	}

	keys := keyboardMapping[player][button]
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}

	return false
}
