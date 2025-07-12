// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

import "math"

// Note represents a musical note as a chromatic semitone index.
//
// 0 corresponds to C-0, 1 to C#0, 2 to D-0, etc.
// Every 12 steps advance one octave.
//
// Typical range for tracker modules (e.g. XM) is:
//
//	0 = C-0
//	1 = C#0
//	...
//	12 = C-1
//	...
//	95 = B-7
//
// In general:
//
//	Note = Octave * 12 + Semitone
type Note = uint8

// NoteToFreq converts a tracker-style Note number to its frequency in Hz.
//
// Uses standard 12-tone equal temperament tuning with A4 = 440 Hz.
//
// Example:
//
//	NoteToFreq(69) // => 440.0 (A4)
//	NoteToFreq(60) // => ~261.63 (C4)
//
// Formula:
//
//	f = 440 * 2^((note - 69) / 12)
func NoteToFreq(note Note) Freq {
	const A4Note = 69
	const A4Freq = 440.0
	n := float64(note - A4Note)
	return A4Freq * math.Pow(2.0, n/12.0)
}
