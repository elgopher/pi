// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pikey provides functionality for a virtual keyboard.
//
// It has the following layout:
//
//	Esc   F1 F2 F3 F4 F5 F6 F7 F8 F9 F10 F11 F12
//	`      1  2  3  4  5  6  7  8  9  0  -  =  Backspace
//	Tab     Q  W  E  R  T  Y  U  I  O  P  [  ]  \
//	CapsLock  A  S  D  F  G  H  J  K  L  ;  '  Enter
//	ShiftLeft  Z  X  C  V  B  N  M   ,  .  /    ShiftRight          Up
//	CtrlLeft AltLeft                    AltRight  CtrlRight   Left Down Right
//
// The pikey package supports up to 119 keys but defines only
// 75 public constants. This is intentional, encouraging you to design your game
// to use no more than 75 buttons. Why? Because this ensures compatibility
// with nearly any modern keyboard. Of course, you can allow
// users to map additional keys if they want — for example,
// by capturing the key name from an Event and saving
// it in a settings file. Just remember that some keyboards
// have even more than 119 keys, but Pi does not support them.
package pikey

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/internal/input"
	"github.com/elgopher/pi/pievent"
)

// Duration returns the number of frames the key k has been held down.
func Duration(k Key) int {
	return keyState.Duration(k)
}

type Key string

const (
	A            Key = "A"
	B            Key = "B"
	C            Key = "C"
	D            Key = "D"
	E            Key = "E"
	F            Key = "F"
	G            Key = "G"
	H            Key = "H"
	I            Key = "I"
	J            Key = "J"
	K            Key = "K"
	L            Key = "L"
	M            Key = "M"
	N            Key = "N"
	O            Key = "O"
	P            Key = "P"
	Q            Key = "Q"
	R            Key = "R"
	S            Key = "S"
	T            Key = "T"
	U            Key = "U"
	V            Key = "V"
	W            Key = "W"
	X            Key = "X"
	Y            Key = "Y"
	Z            Key = "Z"
	AltLeft      Key = "AltLeft"
	AltRight     Key = "AltRight"
	Down         Key = "Down"
	Left         Key = "Left"
	Right        Key = "Right"
	Up           Key = "Up"
	Backquote    Key = "`"
	Backslash    Key = `\`
	Backspace    Key = "Backspace"
	BracketLeft  Key = "BracketLeft"
	BracketRight Key = "BracketRight"
	CapsLock     Key = "CapsLock"
	Comma        Key = ","
	CtrlLeft     Key = "CtrlLeft"
	CtrlRight    Key = "CtrlRight"
	Digit0       Key = "0"
	Digit1       Key = "1"
	Digit2       Key = "2"
	Digit3       Key = "3"
	Digit4       Key = "4"
	Digit5       Key = "5"
	Digit6       Key = "6"
	Digit7       Key = "7"
	Digit8       Key = "8"
	Digit9       Key = "9"
	Enter        Key = "Enter"
	Equal        Key = "="
	Esc          Key = "Esc"
	F1           Key = "F1"
	F2           Key = "F2"
	F3           Key = "F3"
	F4           Key = "F4"
	F5           Key = "F5"
	F6           Key = "F6"
	F7           Key = "F7"
	F8           Key = "F8"
	F9           Key = "F9"
	F10          Key = "F10"
	F11          Key = "F11"
	F12          Key = "F12"
	Minus        Key = "Minus"
	Period       Key = "."
	Quote        Key = `'`
	Semicolon    Key = ";"
	ShiftLeft    Key = "ShiftLeft"
	ShiftRight   Key = "ShiftRight"
	Slash        Key = "/"
	Space        Key = " "
	Tab          Key = "Tab"
	Alt          Key = "Alt"     // alt left or right
	Control      Key = "Control" // control left or right
	Shift        Key = "Shift"   // shift left or right
)

var target = pievent.NewTarget[Event]()
var debugTarget = pievent.NewTarget[Event]()

var keyState input.State[Key]

func init() {
	onKey := func(event Event, _ pievent.Handler) {
		switch event.Type {
		case EventDown:
			keyState.SetDownFrame(event.Key, pi.Frame)
		case EventUp:
			keyState.SetUpFrame(event.Key, pi.Frame)
		}
	}
	Target().SubscribeAll(onKey)
	DebugTarget().SubscribeAll(onKey)
}
