// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Example of using the high-level audio API in Pi.
//
// This program plays a sample at different pitches depending on the key pressed.
package main

import (
	_ "embed"
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pikey"
	"math"
	"slices"
)

var (
	//go:embed "piano.raw"
	sampleRAW []byte
	//go:embed "piano.png"
	pianoPNG []byte
)

func main() {
	// Decode a raw sample file (8-bit mono, no header, no compression).
	sample := piaudio.DecodeRaw(sampleRAW, 16726)

	var selectedKey pikey.Key

	// Load palette and canvas from PNG file
	pi.Palette = pi.DecodePalette(pianoPNG)
	pianoCanvas := pi.DecodeCanvas(pianoPNG)

	pi.SetScreenSize(113, 46)
	pi.Palt(0, false) // disable transparency for color 0

	pi.Init = func() {
		// The sample must be loaded before use,
		// but communication with the audio backend is only possible after starting the game.
		piaudio.LoadSample(sample)
	}

	pi.Update = func() {
		// Check if any of the piano keys were just pressed
		for key, pitch := range buttonPitch {
			if pikey.Duration(key) == 1 {
				selectedKey = key
				// Play the sample on two channels: left and right — for stereo effect
				ch := piaudio.Chan1 | piaudio.Chan2
				vol := 1.0
				piaudio.Play(ch, sample, pitch, vol)
				break
			}
		}
	}

	pi.Draw = func() {
		pi.Cls()

		// Map each key color to either white (tone) or black (semitone)
		for key, color := range keyColors {
			if slices.Contains(semitoneLetters, key) {
				pi.Pal(color, 0)
			} else {
				pi.Pal(color, 7)
			}
		}

		// Highlight the last pressed key in gray
		if selectedKey != "" {
			pi.Pal(keyColors[selectedKey], 6)
		}

		// Draw the piano image with updated color tables
		pi.Blit(pianoCanvas, 0, 0)

		// Draw labels for each piano key
		printLetters()
	}

	piebiten.Run()
}

// cPitch = 1.0 is the base pitch (e.g. middle C).
// Change to 2.0 to play one octave higher, or 0.5 for one octave lower.
const cPitch = 1.0

// Maps keyboard keys to pitch multipliers.
var buttonPitch = map[pikey.Key]float64{
	pikey.Z:     cPitch,          // C
	pikey.S:     adjustPitch(1),  // C#
	pikey.X:     adjustPitch(2),  // D
	pikey.D:     adjustPitch(3),  // D#
	pikey.C:     adjustPitch(4),  // E
	pikey.V:     adjustPitch(5),  // F
	pikey.G:     adjustPitch(6),  // F#
	pikey.B:     adjustPitch(7),  // G
	pikey.H:     adjustPitch(8),  // G#
	pikey.N:     adjustPitch(9),  // A
	pikey.J:     adjustPitch(10), // A#
	pikey.M:     adjustPitch(11), // H
	pikey.Comma: adjustPitch(12), // C
}

func adjustPitch(i int) float64 {
	return cPitch * math.Pow(2, float64(i)/12.0)
}

// Layout of tone and semitone key labels on the piano image.
var (
	toneLetters     = []pikey.Key{"Z", "X", "C", "V", "B", "N", "M", ","}
	semitoneLetters = []pikey.Key{"S", "D", "", "G", "H", "J"}
)

const keyWidth = 13

func printLetters() {
	pi.SetColor(16)
	for i, letter := range toneLetters {
		picofont.Print(string(letter), 9+i*keyWidth, 31)
	}

	for i, letter := range semitoneLetters {
		picofont.Print(string(letter), 16+i*keyWidth, 13)
	}
}

// Each key on the image uses a unique color
var keyColors = map[pikey.Key]pi.Color{
	pikey.Z:     1,
	pikey.S:     2,
	pikey.X:     3,
	pikey.D:     4,
	pikey.C:     6,
	pikey.V:     8,
	pikey.G:     9,
	pikey.B:     10,
	pikey.H:     11,
	pikey.N:     12,
	pikey.J:     13,
	pikey.M:     14,
	pikey.Comma: 15,
}
