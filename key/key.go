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
)

type Button int

const (
	Shift        = 1
	Ctrl         = 3
	Alt          = 5 // Please note that on some keyboard layouts on Windows the right alt is a combination of Ctrl+Alt
	Cap          = 7
	Back         = '\b' // 8
	Tab          = '\t' // 9
	Enter        = '\n' // 10
	F1           = 11
	F2           = 12
	F3           = 13
	F4           = 14
	F5           = 15
	F6           = 16
	F7           = 17
	F8           = 18
	F9           = 19
	F10          = 20
	F11          = 21
	F12          = 22
	Left         = 23
	Right        = 24
	Up           = 25
	Down         = 26
	Esc          = 27
	Space        = ' '  // 32
	Apostrophe   = '\'' // 39
	Comma        = ','  // 44
	Minus        = '-'  // 45
	Period       = '.'  // 46
	Slash        = '/'  // 47
	Digit0       = '0'  // 48
	Digit1       = '1'  // 49
	Digit2       = '2'  // 50
	Digit3       = '3'  // 51
	Digit4       = '4'  // 52
	Digit5       = '5'  // 53
	Digit6       = '6'  // 54
	Digit7       = '7'  // 55
	Digit8       = '8'  // 56
	Digit9       = '9'  // 57
	Semicolon    = ';'  // 59
	Equal        = '='  // 61
	A            = 'A'  // 65
	B            = 'B'  // 66
	C            = 'C'  // 67
	D            = 'D'  // 68
	E            = 'E'  // 69
	F            = 'F'  // 70
	G            = 'G'  // 71
	H            = 'H'  // 72
	I            = 'I'  // 73
	J            = 'J'  // 74
	K            = 'K'  // 75
	L            = 'L'  // 76
	M            = 'M'  // 77
	N            = 'N'  // 78
	O            = 'O'  // 79
	P            = 'P'  // 80
	Q            = 'Q'  // 81
	R            = 'R'  // 82
	S            = 'S'  // 83
	T            = 'T'  // 84
	U            = 'U'  // 85
	V            = 'V'  // 86
	W            = 'W'  // 87
	X            = 'X'  // 88
	Y            = 'Y'  // 89
	Z            = 'Z'  // 90
	BracketLeft  = '['  // 91
	Backslash    = '\\' // 92
	BracketRight = ']'  // 93
	Backquote    = '`'  // 96
)

// Duration has info how many frames in a row a keyboard button was pressed:
// Index of array is equal to key button constant, for example
// Duration[1] has button duration for Shift key.
var Duration [97]uint

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
	return Duration[adjustedButton(b)] > 0
}

// Btnp returns true when the keyboard button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
func Btnp(b Button) bool {
	return input.IsPressedRepeatably(Duration[adjustedButton(b)])
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
