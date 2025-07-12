// Copyright 2025 Jacek Olszak
// This code is licensed under the MIT license (see LICENSE for details)

// Example of programming audio using the low-level API.
//
// This API can be used, for example, to implement packages
// capable of playing module formats (MOD, XM, etc.) and sound effects.
package main

import (
	_ "embed"
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/pi/pimouse"
	"log"
)

//go:embed "wave.wav"
var clickWav []byte

func main() {
	sample := piaudio.DecodeWav(clickWav, 2349.32)

	piloop.Target().Subscribe(piloop.EventGameStarted, func(piloop.Event, pievent.Handler) {
		// The sample must be loaded before use,
		// but communication with the audio backend is only possible after starting the game.
		piaudio.LoadSample(sample)
	})

	// function that schedules playing the SFX
	scheduleSFX := func() {
		// all commands are scheduled with a minimum delay of 0 seconds.
		// However, this doesn't mean the sound plays instantly.
		// On desktop, the delay is around 20 ms; in browsers, about 60 ms.
		// The backend automatically delays commands to reduce audio glitches.
		delay := 0.0

		// remove all planned commands from the channel
		piaudio.ClearChan(piaudio.Chan1, delay)

		// set the sample to play from the beginning (offset=0)
		piaudio.SetSample(piaudio.Chan1, sample, 0, delay)

		// the sound is very short, so we need to loop it.
		// the loop covers the entire sample.
		piaudio.SetLoop(piaudio.Chan1, 0, sample.Len()-1, piaudio.LoopForward, delay)

		for i := 1.0; i > -0.01; i -= 0.01 {
			// gradually reduce the volume to 0
			piaudio.SetVolume(piaudio.Chan1, i, delay)

			// gradually reduce the pitch down to ~340
			pitch := 540 - delay*200
			piaudio.SetPitch(piaudio.Chan1, pitch, delay)
			delay += 0.01
		}
	}

	leftDown := pimouse.EventButton{
		Type:   pimouse.EventButtonDown,
		Button: pimouse.Left,
	}
	pimouse.ButtonTarget().Subscribe(leftDown, func(pimouse.EventButton, pievent.Handler) {
		scheduleSFX()
	})

	piloop.Target().Subscribe(piloop.EventUpdate, func(piloop.Event, pievent.Handler) {
		// piaudio.Time should be used by code that wants to, for example,
		// record when track playback started.
		log.Println("TIME", piaudio.Time)
	})

	piebiten.Run()
}
