// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/vm"
)

var keyMapping = map[int]ebiten.Key{
	vm.KeyShift:        ebiten.KeyShift,
	vm.KeyCtrl:         ebiten.KeyControl,
	vm.KeyAlt:          ebiten.KeyAlt,
	vm.KeyCap:          ebiten.KeyCapsLock,
	vm.KeyBack:         ebiten.KeyBackspace,
	vm.KeyTab:          ebiten.KeyTab,
	vm.KeyEnter:        ebiten.KeyEnter,
	vm.KeyF1:           ebiten.KeyF1,
	vm.KeyF2:           ebiten.KeyF2,
	vm.KeyF3:           ebiten.KeyF3,
	vm.KeyF4:           ebiten.KeyF4,
	vm.KeyF5:           ebiten.KeyF5,
	vm.KeyF6:           ebiten.KeyF6,
	vm.KeyF7:           ebiten.KeyF7,
	vm.KeyF8:           ebiten.KeyF8,
	vm.KeyF9:           ebiten.KeyF9,
	vm.KeyF10:          ebiten.KeyF10,
	vm.KeyF11:          ebiten.KeyF11,
	vm.KeyF12:          ebiten.KeyF12,
	vm.KeyLeft:         ebiten.KeyArrowLeft,
	vm.KeyRight:        ebiten.KeyArrowRight,
	vm.KeyUp:           ebiten.KeyArrowUp,
	vm.KeyDown:         ebiten.KeyArrowDown,
	vm.KeyEsc:          ebiten.KeyEscape,
	vm.KeySpace:        ebiten.KeySpace,
	vm.KeyApostrophe:   ebiten.KeyApostrophe,
	vm.KeyComma:        ebiten.KeyComma,
	vm.KeyMinus:        ebiten.KeyMinus,
	vm.KeyPeriod:       ebiten.KeyPeriod,
	vm.KeySlash:        ebiten.KeySlash,
	vm.KeyDigit0:       ebiten.KeyDigit0,
	vm.KeyDigit1:       ebiten.KeyDigit1,
	vm.KeyDigit2:       ebiten.KeyDigit2,
	vm.KeyDigit3:       ebiten.KeyDigit3,
	vm.KeyDigit4:       ebiten.KeyDigit4,
	vm.KeyDigit5:       ebiten.KeyDigit5,
	vm.KeyDigit6:       ebiten.KeyDigit6,
	vm.KeyDigit7:       ebiten.KeyDigit7,
	vm.KeyDigit8:       ebiten.KeyDigit8,
	vm.KeyDigit9:       ebiten.KeyDigit9,
	vm.KeySemicolon:    ebiten.KeySemicolon,
	vm.KeyEqual:        ebiten.KeyEqual,
	vm.KeyA:            ebiten.KeyA,
	vm.KeyB:            ebiten.KeyB,
	vm.KeyC:            ebiten.KeyC,
	vm.KeyD:            ebiten.KeyD,
	vm.KeyE:            ebiten.KeyE,
	vm.KeyF:            ebiten.KeyF,
	vm.KeyG:            ebiten.KeyG,
	vm.KeyH:            ebiten.KeyH,
	vm.KeyI:            ebiten.KeyI,
	vm.KeyJ:            ebiten.KeyJ,
	vm.KeyK:            ebiten.KeyK,
	vm.KeyL:            ebiten.KeyL,
	vm.KeyM:            ebiten.KeyM,
	vm.KeyN:            ebiten.KeyN,
	vm.KeyO:            ebiten.KeyO,
	vm.KeyP:            ebiten.KeyP,
	vm.KeyQ:            ebiten.KeyQ,
	vm.KeyR:            ebiten.KeyR,
	vm.KeyS:            ebiten.KeyS,
	vm.KeyT:            ebiten.KeyT,
	vm.KeyU:            ebiten.KeyU,
	vm.KeyV:            ebiten.KeyV,
	vm.KeyW:            ebiten.KeyW,
	vm.KeyX:            ebiten.KeyX,
	vm.KeyY:            ebiten.KeyY,
	vm.KeyZ:            ebiten.KeyZ,
	vm.KeyBracketLeft:  ebiten.KeyBracketLeft,
	vm.KeyBackslash:    ebiten.KeyBackslash,
	vm.KeyBracketRight: ebiten.KeyBracketRight,
	vm.KeyBackquote:    ebiten.KeyBackquote,
}

func updateKeyDuration() {
	for button, key := range keyMapping {
		if ebiten.IsKeyPressed(key) {
			vm.KeyDuration[button]++
		} else {
			vm.KeyDuration[button] = 0
		}
	}
}
