// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/elgopher/pi/pikey"
)

func (g *Backend) updateKeyboard() {
	g.keys = inpututil.AppendJustPressedKeys(g.keys[:0])
	for _, key := range g.keys {
		event := pikey.Event{Type: pikey.EventDown, Key: keyMap[key]}
		if !*g.Paused {
			pikey.Target().Publish(event)
		}
		pikey.DebugTarget().Publish(event)
	}

	g.keys = inpututil.AppendJustReleasedKeys(g.keys[:0])
	for _, key := range g.keys {
		event := pikey.Event{Type: pikey.EventUp, Key: keyMap[key]}
		if !*g.Paused {
			pikey.Target().Publish(event)
		}
		pikey.DebugTarget().Publish(event)
	}
}

var keyMap = map[ebiten.Key]pikey.Key{
	ebiten.KeyA:              pikey.A,
	ebiten.KeyB:              pikey.B,
	ebiten.KeyC:              pikey.C,
	ebiten.KeyD:              pikey.D,
	ebiten.KeyE:              pikey.E,
	ebiten.KeyF:              pikey.F,
	ebiten.KeyG:              pikey.G,
	ebiten.KeyH:              pikey.H,
	ebiten.KeyI:              pikey.I,
	ebiten.KeyJ:              pikey.J,
	ebiten.KeyK:              pikey.K,
	ebiten.KeyL:              pikey.L,
	ebiten.KeyM:              pikey.M,
	ebiten.KeyN:              pikey.N,
	ebiten.KeyO:              pikey.O,
	ebiten.KeyP:              pikey.P,
	ebiten.KeyQ:              pikey.Q,
	ebiten.KeyR:              pikey.R,
	ebiten.KeyS:              pikey.S,
	ebiten.KeyT:              pikey.T,
	ebiten.KeyU:              pikey.U,
	ebiten.KeyV:              pikey.V,
	ebiten.KeyW:              pikey.W,
	ebiten.KeyX:              pikey.X,
	ebiten.KeyY:              pikey.Y,
	ebiten.KeyZ:              pikey.Z,
	ebiten.KeyAltLeft:        pikey.AltLeft,
	ebiten.KeyAltRight:       pikey.AltRight,
	ebiten.KeyArrowDown:      pikey.Down,
	ebiten.KeyArrowLeft:      pikey.Left,
	ebiten.KeyArrowRight:     pikey.Right,
	ebiten.KeyArrowUp:        pikey.Up,
	ebiten.KeyBackquote:      pikey.Backquote,
	ebiten.KeyBackslash:      pikey.Backslash,
	ebiten.KeyBackspace:      pikey.Backspace,
	ebiten.KeyBracketLeft:    pikey.BracketLeft,
	ebiten.KeyBracketRight:   pikey.BracketRight,
	ebiten.KeyCapsLock:       pikey.CapsLock,
	ebiten.KeyComma:          pikey.Comma,
	ebiten.KeyContextMenu:    "ContextMenu",
	ebiten.KeyControlLeft:    pikey.CtrlLeft,
	ebiten.KeyControlRight:   pikey.CtrlRight,
	ebiten.KeyDelete:         "Del",
	ebiten.KeyDigit0:         pikey.Digit0,
	ebiten.KeyDigit1:         pikey.Digit1,
	ebiten.KeyDigit2:         pikey.Digit2,
	ebiten.KeyDigit3:         pikey.Digit3,
	ebiten.KeyDigit4:         pikey.Digit4,
	ebiten.KeyDigit5:         pikey.Digit5,
	ebiten.KeyDigit6:         pikey.Digit6,
	ebiten.KeyDigit7:         pikey.Digit7,
	ebiten.KeyDigit8:         pikey.Digit8,
	ebiten.KeyDigit9:         pikey.Digit9,
	ebiten.KeyEnd:            "End",
	ebiten.KeyEnter:          pikey.Enter,
	ebiten.KeyEqual:          pikey.Equal,
	ebiten.KeyEscape:         pikey.Esc,
	ebiten.KeyF1:             pikey.F1,
	ebiten.KeyF2:             pikey.F2,
	ebiten.KeyF3:             pikey.F3,
	ebiten.KeyF4:             pikey.F4,
	ebiten.KeyF5:             pikey.F5,
	ebiten.KeyF6:             pikey.F6,
	ebiten.KeyF7:             pikey.F7,
	ebiten.KeyF8:             pikey.F8,
	ebiten.KeyF9:             pikey.F9,
	ebiten.KeyF10:            pikey.F10,
	ebiten.KeyF11:            pikey.F11,
	ebiten.KeyF12:            pikey.F12,
	ebiten.KeyF13:            "F13",
	ebiten.KeyF14:            "F14",
	ebiten.KeyF15:            "F15",
	ebiten.KeyF16:            "F16",
	ebiten.KeyF17:            "F17",
	ebiten.KeyF18:            "F18",
	ebiten.KeyF19:            "F19",
	ebiten.KeyF20:            "F20",
	ebiten.KeyF21:            "F21",
	ebiten.KeyF22:            "F22",
	ebiten.KeyF23:            "F23",
	ebiten.KeyF24:            "F24",
	ebiten.KeyHome:           "Home",
	ebiten.KeyInsert:         "Ins",
	ebiten.KeyIntlBackslash:  `Intl \`,
	ebiten.KeyMetaLeft:       "MetaLeft",
	ebiten.KeyMetaRight:      "MetaRight",
	ebiten.KeyMinus:          pikey.Minus,
	ebiten.KeyNumLock:        "NumLock",
	ebiten.KeyNumpad0:        "Num 0",
	ebiten.KeyNumpad1:        "Num 1",
	ebiten.KeyNumpad2:        "Num 2",
	ebiten.KeyNumpad3:        "Num 3",
	ebiten.KeyNumpad4:        "Num 4",
	ebiten.KeyNumpad5:        "Num 5",
	ebiten.KeyNumpad6:        "Num 6",
	ebiten.KeyNumpad7:        "Num 7",
	ebiten.KeyNumpad8:        "Num 8",
	ebiten.KeyNumpad9:        "Num 9",
	ebiten.KeyNumpadAdd:      "Num +",
	ebiten.KeyNumpadDecimal:  "Num .",
	ebiten.KeyNumpadDivide:   "Num /",
	ebiten.KeyNumpadEnter:    "Num Enter",
	ebiten.KeyNumpadEqual:    "Num =",
	ebiten.KeyNumpadMultiply: "Num *",
	ebiten.KeyNumpadSubtract: "Num /",
	ebiten.KeyPageDown:       "PgDown",
	ebiten.KeyPageUp:         "PgUp",
	ebiten.KeyPause:          "Pause",
	ebiten.KeyPeriod:         pikey.Period,
	ebiten.KeyPrintScreen:    "PrintScreen",
	ebiten.KeyQuote:          pikey.Quote,
	ebiten.KeyScrollLock:     "ScrollLock",
	ebiten.KeySemicolon:      pikey.Semicolon,
	ebiten.KeyShiftLeft:      pikey.ShiftLeft,
	ebiten.KeyShiftRight:     pikey.ShiftRight,
	ebiten.KeySlash:          pikey.Slash,
	ebiten.KeySpace:          pikey.Space,
	ebiten.KeyTab:            pikey.Tab,
	ebiten.KeyAlt:            pikey.Alt,
	ebiten.KeyControl:        pikey.Ctrl,
	ebiten.KeyShift:          pikey.Shift,
	ebiten.KeyMeta:           "Meta",
}
