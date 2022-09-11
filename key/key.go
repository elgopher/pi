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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/internal/gameloop"
	"github.com/elgopher/pi/internal/input"
)

type Button int

const (
	Shift        Button = 1
	Ctrl         Button = 3
	Alt          Button = 5 // Please note that on some keyboard layouts on Windows the right alt is a combination of Ctrl+Alt
	Cap          Button = 7
	Back         Button = '\b' // 8
	Tab          Button = '\t' // 9
	Enter        Button = '\n' // 10
	F1           Button = 11
	F2           Button = 12
	F3           Button = 13
	F4           Button = 14
	F5           Button = 15
	F6           Button = 16
	F7           Button = 17
	F8           Button = 18
	F9           Button = 19
	F10          Button = 20
	F11          Button = 21
	F12          Button = 22
	Left         Button = 23
	Right        Button = 24
	Up           Button = 25
	Down         Button = 26
	Esc          Button = 27
	Space        Button = ' '  // 32
	Apostrophe   Button = '\'' // 39
	Comma        Button = ','  // 44
	Minus        Button = '-'  // 45
	Period       Button = '.'  // 46
	Slash        Button = '/'  // 47
	Digit0       Button = '0'  // 48
	Digit1       Button = '1'  // 49
	Digit2       Button = '2'  // 50
	Digit3       Button = '3'  // 51
	Digit4       Button = '4'  // 52
	Digit5       Button = '5'  // 53
	Digit6       Button = '6'  // 54
	Digit7       Button = '7'  // 55
	Digit8       Button = '8'  // 56
	Digit9       Button = '9'  // 57
	Semicolon    Button = ';'  // 59
	Equal        Button = '='  // 61
	A            Button = 'A'  // 65
	B            Button = 'B'  // 66
	C            Button = 'C'  // 67
	D            Button = 'D'  // 68
	E            Button = 'E'  // 69
	F            Button = 'F'  // 70
	G            Button = 'G'  // 71
	H            Button = 'H'  // 72
	I            Button = 'I'  // 73
	J            Button = 'J'  // 74
	K            Button = 'K'  // 75
	L            Button = 'L'  // 76
	M            Button = 'M'  // 77
	N            Button = 'N'  // 78
	O            Button = 'O'  // 79
	P            Button = 'P'  // 80
	Q            Button = 'Q'  // 81
	R            Button = 'R'  // 82
	S            Button = 'S'  // 83
	T            Button = 'T'  // 84
	U            Button = 'U'  // 85
	V            Button = 'V'  // 86
	W            Button = 'W'  // 87
	X            Button = 'X'  // 88
	Y            Button = 'Y'  // 89
	Z            Button = 'Z'  // 90
	BracketLeft  Button = '['  // 91
	Backslash    Button = '\\' // 92
	BracketRight Button = ']'  // 93
	Backquote    Button = '`'  // 96
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
	return ebiten.IsKeyPressed(piToEbitengineMapping[adjustedButton(b)])
}

// Btnp returns true when the keyboard button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
func Btnp(b Button) bool {
	return input.IsPressedRepeatably(buttonDurations[adjustedButton(b)])
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

var buttonDurations [97]int

var piToEbitengineMapping = map[Button]ebiten.Key{
	Shift:        ebiten.KeyShift,
	Ctrl:         ebiten.KeyControl,
	Alt:          ebiten.KeyAlt,
	Cap:          ebiten.KeyCapsLock,
	Back:         ebiten.KeyBackspace,
	Tab:          ebiten.KeyTab,
	Enter:        ebiten.KeyEnter,
	F1:           ebiten.KeyF1,
	F2:           ebiten.KeyF2,
	F3:           ebiten.KeyF3,
	F4:           ebiten.KeyF4,
	F5:           ebiten.KeyF5,
	F6:           ebiten.KeyF6,
	F7:           ebiten.KeyF7,
	F8:           ebiten.KeyF8,
	F9:           ebiten.KeyF9,
	F10:          ebiten.KeyF10,
	F11:          ebiten.KeyF11,
	F12:          ebiten.KeyF12,
	Left:         ebiten.KeyArrowLeft,
	Right:        ebiten.KeyArrowRight,
	Up:           ebiten.KeyArrowUp,
	Down:         ebiten.KeyArrowDown,
	Esc:          ebiten.KeyEscape,
	Space:        ebiten.KeySpace,
	Apostrophe:   ebiten.KeyApostrophe,
	Comma:        ebiten.KeyComma,
	Minus:        ebiten.KeyMinus,
	Period:       ebiten.KeyPeriod,
	Slash:        ebiten.KeySlash,
	Digit0:       ebiten.KeyDigit0,
	Digit1:       ebiten.KeyDigit1,
	Digit2:       ebiten.KeyDigit2,
	Digit3:       ebiten.KeyDigit3,
	Digit4:       ebiten.KeyDigit4,
	Digit5:       ebiten.KeyDigit5,
	Digit6:       ebiten.KeyDigit6,
	Digit7:       ebiten.KeyDigit7,
	Digit8:       ebiten.KeyDigit8,
	Digit9:       ebiten.KeyDigit9,
	Semicolon:    ebiten.KeySemicolon,
	Equal:        ebiten.KeyEqual,
	A:            ebiten.KeyA,
	B:            ebiten.KeyB,
	C:            ebiten.KeyC,
	D:            ebiten.KeyD,
	E:            ebiten.KeyE,
	F:            ebiten.KeyF,
	G:            ebiten.KeyG,
	H:            ebiten.KeyH,
	I:            ebiten.KeyI,
	J:            ebiten.KeyJ,
	K:            ebiten.KeyK,
	L:            ebiten.KeyL,
	M:            ebiten.KeyM,
	N:            ebiten.KeyN,
	O:            ebiten.KeyO,
	P:            ebiten.KeyP,
	Q:            ebiten.KeyQ,
	R:            ebiten.KeyR,
	S:            ebiten.KeyS,
	T:            ebiten.KeyT,
	U:            ebiten.KeyU,
	V:            ebiten.KeyV,
	W:            ebiten.KeyW,
	X:            ebiten.KeyX,
	Y:            ebiten.KeyY,
	Z:            ebiten.KeyZ,
	BracketLeft:  ebiten.KeyBracketLeft,
	Backslash:    ebiten.KeyBackslash,
	BracketRight: ebiten.KeyBracketRight,
	Backquote:    ebiten.KeyBackquote,
}

func init() {
	gameloop.UpdateFunctions = append(gameloop.UpdateFunctions, update)
}

func update() {
	for button, key := range piToEbitengineMapping {
		if ebiten.IsKeyPressed(key) {
			buttonDurations[button]++
		} else {
			buttonDurations[button] = 0
		}
	}
}
