// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Button is a virtual button on any game controller. The game controller can be a gamepad or a keyboard.
//
// Button is used by Btn, Btnp, BtnPlayer, BtnpPlayer, BtnBits and BtnpBits.
type Button int

// Keyboard mappings:
//
//   player 0: [DPAD] - cursors  [O] - Z C N   [X] - X V M
//   player 1: [DPAD] - SFED     [O] - LSHIFT  [X] - TAB W Q A
//
// First connected gamepad controller is player 0, second player 1 and so on.
// On XBox controller [O] is A and Y, [X] is B and X.
const (
	Left  Button = 0
	Right Button = 1
	Up    Button = 2
	Down  Button = 3
	O     Button = 4 // O is a first fire button
	X     Button = 5 // X is a second fire button
)

// Btn returns true if a button is being pressed at this moment by player 0.
func Btn(button Button) bool {
	return BtnPlayer(button, 0)
}

// BtnPlayer returns true if a button is being pressed at this moment
// by specific player. The player can be 0..7.
func BtnPlayer(button Button, player int) bool {
	return isPressed(buttonDuration(player, button))
}

// Btnp returns true when the button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
// This simulates keyboard-like repeating.
func Btnp(button Button) bool {
	return BtnpPlayer(button, 0)
}

// BtnpPlayer returns true when the button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
// This simulates keyboard-like repeating. The player can be 0..7.
func BtnpPlayer(button Button, player int) bool {
	return isPressedRepeatably(buttonDuration(player, button))
}

// BtnBits returns the state of all buttons for players 0 and 1 as bitset.
//
// The first byte contains the button states for player 0 (bits 0 through 5, bits 6 and 7 are unused).
// The second byte contains the button states for player 1 (bits 8 through 13).
//
// Bit 0 is Left, 1 is Right, bit 5 is the X button.
//
// A bit of 1 means the button is pressed.
func BtnBits() int {
	return controllers[0].bits(isPressed) + controllers[1].bits(isPressed)<<8
}

// BtnpBits returns the state of all buttons for players 0 and 1 as bitset.
//
// The first byte contains the button states for player 0 (bits 0 through 5, bits 6 and 7 are unused).
// The second byte contains the button states for player 1 (bits 8 through 13).
//
// Bit 0 is Left, 1 is Right, bit 5 is the X button.
//
// A bit of 1 means the button has just been pressed.
func BtnpBits() int {
	return controllers[0].bits(isPressedRepeatably) + controllers[1].bits(isPressedRepeatably)<<8
}

func buttonDuration(player int, button Button) int {
	if button < Left || button > X {
		return 0
	}

	if player < 0 || player > 7 {
		return 0
	}

	return controllers[player].buttonDuration[button]
}

func isPressed(duration int) bool {
	return duration > 0
}

func isPressedRepeatably(duration int) bool {
	const (
		pressDuration = 15 // make it configurable
		pressInterval = 4  // make it configurable
	)

	if duration == 1 {
		return true
	}

	return duration >= pressDuration+1 && duration%pressInterval == 0
}

func updateController() {
	for player := 0; player < 8; player++ {
		controllers[player].update(player)
	}
}

var controllers [8]controller

type controller struct {
	buttonDuration [6]int // 	left, right, up, down, o, x
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
		isKeyboardPressed(player, Left) {
		c.buttonDuration[Left] += 1
		c.buttonDuration[Right] = 0
	} else if axisX > 0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftRight) ||
		isKeyboardPressed(player, Right) {
		c.buttonDuration[Right] += 1
		c.buttonDuration[Left] = 0
	} else {
		c.buttonDuration[Right] = 0
		c.buttonDuration[Left] = 0
	}

	if axisY < -0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftTop) ||
		isKeyboardPressed(player, Up) {
		c.buttonDuration[Up] += 1
		c.buttonDuration[Down] = 0
	} else if axisY > 0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftBottom) ||
		isKeyboardPressed(player, Down) {
		c.buttonDuration[Down] += 1
		c.buttonDuration[Up] = 0
	} else {
		c.buttonDuration[Up] = 0
		c.buttonDuration[Down] = 0
	}
}

func (c *controller) updateFireButtons(player int) {
	gamepadID := ebiten.GamepadID(player)

	if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightBottom) ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightTop) ||
		isKeyboardPressed(player, O) {
		c.buttonDuration[O] += 1
	} else {
		c.buttonDuration[O] = 0
	}

	if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightRight) ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightLeft) ||
		isKeyboardPressed(player, X) {
		c.buttonDuration[X] += 1
	} else {
		c.buttonDuration[X] = 0
	}
}

func (c *controller) bits(isSet func(int) bool) int {
	var b int
	for i := 0; i <= int(X); i++ {
		if isSet(c.buttonDuration[i]) {
			b += 1 << i
		}
	}
	return b
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

func isKeyboardPressed(player int, button Button) bool {
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
