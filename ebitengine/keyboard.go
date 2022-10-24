// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/mem"
)

var keyMapping = map[int]ebiten.Key{
	mem.KeyShift:        ebiten.KeyShift,
	mem.KeyCtrl:         ebiten.KeyControl,
	mem.KeyAlt:          ebiten.KeyAlt,
	mem.KeyCap:          ebiten.KeyCapsLock,
	mem.KeyBack:         ebiten.KeyBackspace,
	mem.KeyTab:          ebiten.KeyTab,
	mem.KeyEnter:        ebiten.KeyEnter,
	mem.KeyF1:           ebiten.KeyF1,
	mem.KeyF2:           ebiten.KeyF2,
	mem.KeyF3:           ebiten.KeyF3,
	mem.KeyF4:           ebiten.KeyF4,
	mem.KeyF5:           ebiten.KeyF5,
	mem.KeyF6:           ebiten.KeyF6,
	mem.KeyF7:           ebiten.KeyF7,
	mem.KeyF8:           ebiten.KeyF8,
	mem.KeyF9:           ebiten.KeyF9,
	mem.KeyF10:          ebiten.KeyF10,
	mem.KeyF11:          ebiten.KeyF11,
	mem.KeyF12:          ebiten.KeyF12,
	mem.KeyLeft:         ebiten.KeyArrowLeft,
	mem.KeyRight:        ebiten.KeyArrowRight,
	mem.KeyUp:           ebiten.KeyArrowUp,
	mem.KeyDown:         ebiten.KeyArrowDown,
	mem.KeyEsc:          ebiten.KeyEscape,
	mem.KeySpace:        ebiten.KeySpace,
	mem.KeyApostrophe:   ebiten.KeyApostrophe,
	mem.KeyComma:        ebiten.KeyComma,
	mem.KeyMinus:        ebiten.KeyMinus,
	mem.KeyPeriod:       ebiten.KeyPeriod,
	mem.KeySlash:        ebiten.KeySlash,
	mem.KeyDigit0:       ebiten.KeyDigit0,
	mem.KeyDigit1:       ebiten.KeyDigit1,
	mem.KeyDigit2:       ebiten.KeyDigit2,
	mem.KeyDigit3:       ebiten.KeyDigit3,
	mem.KeyDigit4:       ebiten.KeyDigit4,
	mem.KeyDigit5:       ebiten.KeyDigit5,
	mem.KeyDigit6:       ebiten.KeyDigit6,
	mem.KeyDigit7:       ebiten.KeyDigit7,
	mem.KeyDigit8:       ebiten.KeyDigit8,
	mem.KeyDigit9:       ebiten.KeyDigit9,
	mem.KeySemicolon:    ebiten.KeySemicolon,
	mem.KeyEqual:        ebiten.KeyEqual,
	mem.KeyA:            ebiten.KeyA,
	mem.KeyB:            ebiten.KeyB,
	mem.KeyC:            ebiten.KeyC,
	mem.KeyD:            ebiten.KeyD,
	mem.KeyE:            ebiten.KeyE,
	mem.KeyF:            ebiten.KeyF,
	mem.KeyG:            ebiten.KeyG,
	mem.KeyH:            ebiten.KeyH,
	mem.KeyI:            ebiten.KeyI,
	mem.KeyJ:            ebiten.KeyJ,
	mem.KeyK:            ebiten.KeyK,
	mem.KeyL:            ebiten.KeyL,
	mem.KeyM:            ebiten.KeyM,
	mem.KeyN:            ebiten.KeyN,
	mem.KeyO:            ebiten.KeyO,
	mem.KeyP:            ebiten.KeyP,
	mem.KeyQ:            ebiten.KeyQ,
	mem.KeyR:            ebiten.KeyR,
	mem.KeyS:            ebiten.KeyS,
	mem.KeyT:            ebiten.KeyT,
	mem.KeyU:            ebiten.KeyU,
	mem.KeyV:            ebiten.KeyV,
	mem.KeyW:            ebiten.KeyW,
	mem.KeyX:            ebiten.KeyX,
	mem.KeyY:            ebiten.KeyY,
	mem.KeyZ:            ebiten.KeyZ,
	mem.KeyBracketLeft:  ebiten.KeyBracketLeft,
	mem.KeyBackslash:    ebiten.KeyBackslash,
	mem.KeyBracketRight: ebiten.KeyBracketRight,
	mem.KeyBackquote:    ebiten.KeyBackquote,
}

func updateKeyDuration() {
	for button, key := range keyMapping {
		if ebiten.IsKeyPressed(key) {
			mem.KeyDuration[button]++
		} else {
			mem.KeyDuration[button] = 0
		}
	}
}
