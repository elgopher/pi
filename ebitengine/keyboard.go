// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/key"
)

var keyMapping = map[int]ebiten.Key{
	key.Shift:        ebiten.KeyShift,
	key.Ctrl:         ebiten.KeyControl,
	key.Alt:          ebiten.KeyAlt,
	key.Cap:          ebiten.KeyCapsLock,
	key.Back:         ebiten.KeyBackspace,
	key.Tab:          ebiten.KeyTab,
	key.Enter:        ebiten.KeyEnter,
	key.F1:           ebiten.KeyF1,
	key.F2:           ebiten.KeyF2,
	key.F3:           ebiten.KeyF3,
	key.F4:           ebiten.KeyF4,
	key.F5:           ebiten.KeyF5,
	key.F6:           ebiten.KeyF6,
	key.F7:           ebiten.KeyF7,
	key.F8:           ebiten.KeyF8,
	key.F9:           ebiten.KeyF9,
	key.F10:          ebiten.KeyF10,
	key.F11:          ebiten.KeyF11,
	key.F12:          ebiten.KeyF12,
	key.Left:         ebiten.KeyArrowLeft,
	key.Right:        ebiten.KeyArrowRight,
	key.Up:           ebiten.KeyArrowUp,
	key.Down:         ebiten.KeyArrowDown,
	key.Esc:          ebiten.KeyEscape,
	key.Space:        ebiten.KeySpace,
	key.Apostrophe:   ebiten.KeyApostrophe,
	key.Comma:        ebiten.KeyComma,
	key.Minus:        ebiten.KeyMinus,
	key.Period:       ebiten.KeyPeriod,
	key.Slash:        ebiten.KeySlash,
	key.Digit0:       ebiten.KeyDigit0,
	key.Digit1:       ebiten.KeyDigit1,
	key.Digit2:       ebiten.KeyDigit2,
	key.Digit3:       ebiten.KeyDigit3,
	key.Digit4:       ebiten.KeyDigit4,
	key.Digit5:       ebiten.KeyDigit5,
	key.Digit6:       ebiten.KeyDigit6,
	key.Digit7:       ebiten.KeyDigit7,
	key.Digit8:       ebiten.KeyDigit8,
	key.Digit9:       ebiten.KeyDigit9,
	key.Semicolon:    ebiten.KeySemicolon,
	key.Equal:        ebiten.KeyEqual,
	key.A:            ebiten.KeyA,
	key.B:            ebiten.KeyB,
	key.C:            ebiten.KeyC,
	key.D:            ebiten.KeyD,
	key.E:            ebiten.KeyE,
	key.F:            ebiten.KeyF,
	key.G:            ebiten.KeyG,
	key.H:            ebiten.KeyH,
	key.I:            ebiten.KeyI,
	key.J:            ebiten.KeyJ,
	key.K:            ebiten.KeyK,
	key.L:            ebiten.KeyL,
	key.M:            ebiten.KeyM,
	key.N:            ebiten.KeyN,
	key.O:            ebiten.KeyO,
	key.P:            ebiten.KeyP,
	key.Q:            ebiten.KeyQ,
	key.R:            ebiten.KeyR,
	key.S:            ebiten.KeyS,
	key.T:            ebiten.KeyT,
	key.U:            ebiten.KeyU,
	key.V:            ebiten.KeyV,
	key.W:            ebiten.KeyW,
	key.X:            ebiten.KeyX,
	key.Y:            ebiten.KeyY,
	key.Z:            ebiten.KeyZ,
	key.BracketLeft:  ebiten.KeyBracketLeft,
	key.Backslash:    ebiten.KeyBackslash,
	key.BracketRight: ebiten.KeyBracketRight,
	key.Backquote:    ebiten.KeyBackquote,
}

func updateKeyDuration() {
	for button, k := range keyMapping {
		if ebiten.IsKeyPressed(k) {
			key.Duration[button]++
		} else {
			key.Duration[button] = 0
		}
	}
}
