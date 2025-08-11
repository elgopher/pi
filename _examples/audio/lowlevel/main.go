// Copyright 2025 Jacek Olszak
// This code is licensed under the MIT license (see LICENSE for details)

// Example of programming audio using the low-level API.
//
// This API can be used, for example, to implement packages
// capable of playing module formats (MOD, XM, etc.) and sound effects.
package main

import (
	_ "embed"
	"log"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/pi/pimouse"
)

//go:embed "wave.wav"
var clickWav []byte

func main() {
	sample := piaudio.DecodeWav(clickWav)

	pi.Init = func() {
		// The sample must be loaded before use,
		// but communication with the audio backend is only possible after starting the game.
		piaudio.LoadSample(sample)
	}

	// function that schedules playing the SFX
	scheduleSFX := func() {
		// All commands are scheduled with a minimum delay of 0 seconds.
		// However, this doesn't mean the sound plays instantly.
		// On desktop, the delay is around 20 ms; in browsers, about 60 ms.
		// The backend automatically delays commands to reduce audio glitches.
		delay := 0.0

		// Use two channels at once.
		// Chan1 is sent to a left speaker, Chan2 is sent to a right speaker.
		ch := piaudio.Chan1 | piaudio.Chan2

		// remove all planned commands from channels
		piaudio.ClearChan(ch, delay)

		// set the sample to play from the beginning (offset=0)
		piaudio.SetSample(ch, sample, 0, delay)

		// the sound is very short, so we need to loop it.
		// the loop covers the entire sample.
		piaudio.SetLoop(ch, 0, sample.Len(), piaudio.LoopForward, delay)

		for i := 1.0; i > -0.01; i -= 0.01 {
			// gradually reduce the volume to 0
			piaudio.SetVolume(ch, i, delay)

			// gradually reduce the pitch down to 0
			pitch := 1.0 - delay
			piaudio.SetPitch(ch, pitch, delay)
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
