// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package main

import (
	_ "embed"
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/piloop"
)

//go:embed "click.raw"
var clickRaw []byte

func main() {
	sample := piaudio.NewSample(clickRaw)
	sound := piaudio.Sound{
		Sample:   sample,
		BaseFreq: 11025,
		Loop: piaudio.Loop{
			LoopMode: piaudio.LoopModeNone,
		},
		Pitch:  piaudio.A5,
		Offset: 0,
		Volume: 1.0,
	}

	piloop.Target().Subscribe(piloop.EventGameStarted, func(piloop.Event, pievent.Handler) {
		piaudio.Play(piaudio.Chan0, sound, 0) // play the sound immediately
		piaudio.Stop(piaudio.Chan0, 100)      // stop it after 100 ticks

		piaudio.SetPitch(piaudio.Chan0, piaudio.A5+10, 10) // update channel pitch after 10 ticks
		piaudio.SetPitch(piaudio.Chan0, piaudio.A5+20, 20)
		piaudio.SetPitch(piaudio.Chan0, piaudio.A5+30, 30)
		piaudio.SetVolume(piaudio.Chan0, 0.5, 70) // update channel volume after 70 ticks
		piaudio.SetVolume(piaudio.Chan0, 0.3, 80)
		piaudio.SetVolume(piaudio.Chan0, 0.1, 90)
	})

	piebiten.Run()

}
