// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Example of programming audio using the high-level API.
//
// This example plays a sample with different pitches.
package main

import (
	_ "embed"
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pikey"
	"math"
)

//go:embed "piano.raw"
var sampleRAW []byte

const c = 1.0 // change to 2.0 to octave up, change to 0.5 to octave down.

var buttonPitch = map[pikey.Key]float64{
	pikey.Z:     c,                          // C
	pikey.S:     c * math.Pow(2, 1.0/12.0),  // C#
	pikey.X:     c * math.Pow(2, 2.0/12.0),  // D
	pikey.D:     c * math.Pow(2, 3.0/12.0),  // D#
	pikey.C:     c * math.Pow(2, 4.0/12.0),  // E
	pikey.V:     c * math.Pow(2, 5.0/12.0),  // F
	pikey.G:     c * math.Pow(2, 6.0/12.0),  // F#
	pikey.B:     c * math.Pow(2, 7.0/12.0),  // G
	pikey.H:     c * math.Pow(2, 8.0/12.0),  // G#
	pikey.N:     c * math.Pow(2, 9.0/12.0),  // A
	pikey.J:     c * math.Pow(2, 10.0/12.0), // A#
	pikey.M:     c * math.Pow(2, 11.0/12.0), // H
	pikey.Comma: c * math.Pow(2, 12.0/12.0), // C
}

//go:embed "piano.png"
var pianoPNG []byte

var buttonRectangles = map[pikey.Key]pi.IntArea{
	pikey.Z: {5, 5, 16, 40},
	pikey.X: {18, 5, 16 + 13, 40},
}

func main() {
	sample := piaudio.DecodeRaw(sampleRAW, 16726)

	var selectedKey pikey.Key

	//pi.Palette = pi.DecodePalette(pianoPNG)
	canvas := pi.DecodeCanvas(pianoPNG)
	pi.SetScreenSize(113, 46)
	pi.Palt(0, false)

	pi.ColorTables[0][63][7] = 6 // drawing 63 on top of 7 will result in 6

	pi.Init = func() {
		piaudio.LoadSample(sample)
	}

	pi.Update = func() {
		for key, pitch := range buttonPitch {
			if pikey.Duration(key) == 1 {
				selectedKey = key
				ch := piaudio.Chan1 | piaudio.Chan2 // stereo
				piaudio.Play(ch, sample, pitch, 1.0)
				return
			}
		}
	}

	pi.Draw = func() {
		pi.Blit(canvas, 0, 0)
		pi.SetColor(1)
		picofont.Print("Z", 9, 31)
		picofont.Print("X", 22, 31)
		picofont.Print("C", 35, 31)
		picofont.Print("V", 48, 31)
		picofont.Print("B", 61, 31)
		picofont.Print("N", 74, 31)
		picofont.Print("M", 87, 31)
		picofont.Print(",", 100, 31)
		//
		pi.SetColor(7)
		picofont.Print("S", 16, 13)
		picofont.Print("D", 29, 13)
		picofont.Print("G", 55, 13)
		picofont.Print("H", 68, 13)
		picofont.Print("J", 81, 13)

		if selectedKey != "" {
			pi.SetColor(63)
			area := buttonRectangles[selectedKey]
			pi.RectFill(area.X, area.Y, area.W, area.H) // Use W and H as a X2, Y2
		}
	}

	piebiten.Run()
}
