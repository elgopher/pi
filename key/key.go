// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package key provides functions for handling virtual keyboard input.
// Useful for writing tools or games using mouse + keyboard combination.
// For basic game control please consider virtual controller instead (pi.Btn and pi.Btnp).
//
// Virtual keyboard is inspired by US keyboard layout:
//
//	ESC F1 F2 F3 F4 F5 F6 F7 F8 F9 F10 F11 F12
//	`  1  2  3  4  5  6  7  8  9  0  -  =  <--
//	TAB Q  W  E  R  T  Y  U  I  O  P  [   ]  \
//	CAP  A  S  D  F  G  H  J  K  L  ;  ' ENTER
//	SHIFT Z  X  C  V  B  N  M  ,  .  /    ↑
//	CTRL ALT        SPACE              ←  ↓  →
//
// This package is experimental. Can be changed in the future.
package key

import (
	"strconv"
	"unicode"

	"github.com/elgopher/pi/internal/input"
	"github.com/elgopher/pi/vm"
)

type Button int

const (
	Shift        Button = vm.KeyShift
	Ctrl         Button = vm.KeyCtrl
	Alt          Button = vm.KeyAlt // Please note that on some keyboard layouts on Windows the right alt is a combination of Ctrl+Alt
	Cap          Button = vm.KeyCap
	Back         Button = vm.KeyBack
	Tab          Button = vm.KeyTab
	Enter        Button = vm.KeyEnter
	F1           Button = vm.KeyF1
	F2           Button = vm.KeyF2
	F3           Button = vm.KeyF3
	F4           Button = vm.KeyF4
	F5           Button = vm.KeyF5
	F6           Button = vm.KeyF6
	F7           Button = vm.KeyF7
	F8           Button = vm.KeyF8
	F9           Button = vm.KeyF9
	F10          Button = vm.KeyF10
	F11          Button = vm.KeyF11
	F12          Button = vm.KeyF12
	Left         Button = vm.KeyLeft
	Right        Button = vm.KeyRight
	Up           Button = vm.KeyUp
	Down         Button = vm.KeyDown
	Esc          Button = vm.KeyEsc
	Space        Button = vm.KeySpace
	Apostrophe   Button = vm.KeyApostrophe
	Comma        Button = vm.KeyComma
	Minus        Button = vm.KeyMinus
	Period       Button = vm.KeyPeriod
	Slash        Button = vm.KeySlash
	Digit0       Button = vm.KeyDigit0
	Digit1       Button = vm.KeyDigit1
	Digit2       Button = vm.KeyDigit2
	Digit3       Button = vm.KeyDigit3
	Digit4       Button = vm.KeyDigit4
	Digit5       Button = vm.KeyDigit5
	Digit6       Button = vm.KeyDigit6
	Digit7       Button = vm.KeyDigit7
	Digit8       Button = vm.KeyDigit8
	Digit9       Button = vm.KeyDigit9
	Semicolon    Button = vm.KeySemicolon
	Equal        Button = vm.KeyEqual
	A            Button = vm.KeyA
	B            Button = vm.KeyB
	C            Button = vm.KeyC
	D            Button = vm.KeyD
	E            Button = vm.KeyE
	F            Button = vm.KeyF
	G            Button = vm.KeyG
	H            Button = vm.KeyH
	I            Button = vm.KeyI
	J            Button = vm.KeyJ
	K            Button = vm.KeyK
	L            Button = vm.KeyL
	M            Button = vm.KeyM
	N            Button = vm.KeyN
	O            Button = vm.KeyO
	P            Button = vm.KeyP
	Q            Button = vm.KeyQ
	R            Button = vm.KeyR
	S            Button = vm.KeyS
	T            Button = vm.KeyT
	U            Button = vm.KeyU
	V            Button = vm.KeyV
	W            Button = vm.KeyW
	X            Button = vm.KeyX
	Y            Button = vm.KeyY
	Z            Button = vm.KeyZ
	BracketLeft  Button = vm.KeyBracketLeft
	Backslash    Button = vm.KeyBackslash
	BracketRight Button = vm.KeyBracketRight
	Backquote    Button = vm.KeyBackquote
)

// Btn returns true if the keyboard button is being pressed at this moment.
//
// You can use button constants defined in the package or pass runes as
// a button parameter, for example:
//
//	key.Btn(key.A)
//	key.Btn('A')
//	key.Btn('a')
//
// All these calls have same effect.
func Btn(b Button) bool {
	return vm.KeyDuration[adjustedButton(b)] > 0
}

// Btnp returns true when the keyboard button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
func Btnp(b Button) bool {
	return input.IsPressedRepeatably(vm.KeyDuration[adjustedButton(b)])
}

var specialChars = map[Button]string{
	Space: "Space",
	Shift: "Shift",
	Ctrl:  "Ctrl",
	Alt:   "Alt",
	Cap:   "Cap",
	Back:  "Back",
	Tab:   "Tab",
	Enter: "Enter",
	Esc:   "Esc",
	Left:  "Left", Up: "Up", Right: "Right", Down: "Down",
}

func (b Button) String() string {
	printable := b > Space && b <= '~'
	if printable {
		return string(rune(b))
	}

	if b >= F1 && b <= F12 {
		number := int(b-F1) + 1
		return "F" + strconv.Itoa(number)
	}

	if special, ok := specialChars[b]; ok {
		return special
	}

	return "?"
}

func adjustedButton(k Button) Button {
	if k >= 'a' && k <= 'z' {
		return Button(unicode.ToUpper(rune(k)))
	}

	return k
}
