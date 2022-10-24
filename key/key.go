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
	"github.com/elgopher/pi/mem"
)

type Button int

const (
	Shift        Button = mem.KeyShift
	Ctrl         Button = mem.KeyCtrl
	Alt          Button = mem.KeyAlt // Please note that on some keyboard layouts on Windows the right alt is a combination of Ctrl+Alt
	Cap          Button = mem.KeyCap
	Back         Button = mem.KeyBack
	Tab          Button = mem.KeyTab
	Enter        Button = mem.KeyEnter
	F1           Button = mem.KeyF1
	F2           Button = mem.KeyF2
	F3           Button = mem.KeyF3
	F4           Button = mem.KeyF4
	F5           Button = mem.KeyF5
	F6           Button = mem.KeyF6
	F7           Button = mem.KeyF7
	F8           Button = mem.KeyF8
	F9           Button = mem.KeyF9
	F10          Button = mem.KeyF10
	F11          Button = mem.KeyF11
	F12          Button = mem.KeyF12
	Left         Button = mem.KeyLeft
	Right        Button = mem.KeyRight
	Up           Button = mem.KeyUp
	Down         Button = mem.KeyDown
	Esc          Button = mem.KeyEsc
	Space        Button = mem.KeySpace
	Apostrophe   Button = mem.KeyApostrophe
	Comma        Button = mem.KeyComma
	Minus        Button = mem.KeyMinus
	Period       Button = mem.KeyPeriod
	Slash        Button = mem.KeySlash
	Digit0       Button = mem.KeyDigit0
	Digit1       Button = mem.KeyDigit1
	Digit2       Button = mem.KeyDigit2
	Digit3       Button = mem.KeyDigit3
	Digit4       Button = mem.KeyDigit4
	Digit5       Button = mem.KeyDigit5
	Digit6       Button = mem.KeyDigit6
	Digit7       Button = mem.KeyDigit7
	Digit8       Button = mem.KeyDigit8
	Digit9       Button = mem.KeyDigit9
	Semicolon    Button = mem.KeySemicolon
	Equal        Button = mem.KeyEqual
	A            Button = mem.KeyA
	B            Button = mem.KeyB
	C            Button = mem.KeyC
	D            Button = mem.KeyD
	E            Button = mem.KeyE
	F            Button = mem.KeyF
	G            Button = mem.KeyG
	H            Button = mem.KeyH
	I            Button = mem.KeyI
	J            Button = mem.KeyJ
	K            Button = mem.KeyK
	L            Button = mem.KeyL
	M            Button = mem.KeyM
	N            Button = mem.KeyN
	O            Button = mem.KeyO
	P            Button = mem.KeyP
	Q            Button = mem.KeyQ
	R            Button = mem.KeyR
	S            Button = mem.KeyS
	T            Button = mem.KeyT
	U            Button = mem.KeyU
	V            Button = mem.KeyV
	W            Button = mem.KeyW
	X            Button = mem.KeyX
	Y            Button = mem.KeyY
	Z            Button = mem.KeyZ
	BracketLeft  Button = mem.KeyBracketLeft
	Backslash    Button = mem.KeyBackslash
	BracketRight Button = mem.KeyBracketRight
	Backquote    Button = mem.KeyBackquote
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
	return mem.KeyDuration[adjustedButton(b)] > 0
}

// Btnp returns true when the keyboard button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
func Btnp(b Button) bool {
	return input.IsPressedRepeatably(mem.KeyDuration[adjustedButton(b)])
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
