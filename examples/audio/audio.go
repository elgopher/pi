package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio"
	"github.com/elgopher/pi/ebitengine"
)

func main() {
	audio.Sfx[0].Notes[0] = audio.Note{
		Pitch:      audio.PitchG3,
		Instrument: audio.InstrumentPhaser,
		Volume:     7,
	}
	audio.Sfx[0].Speed = 255
	audio.Sync()

	sfxNote0 := &audio.Sfx[0].Notes[0]

	pi.Update = func() {
		if pi.Btnp(pi.X) {
			audio.Play(0, 0, 0, 0)
		}

		if pi.Btnp(pi.Up) && sfxNote0.Pitch < 255 {
			sfxNote0.Pitch += 1
			playSfx0()
		}
		if pi.Btnp(pi.Down) && sfxNote0.Pitch > 0 {
			sfxNote0.Pitch -= 1
			playSfx0()
		}
		if pi.Btnp(pi.Right) && sfxNote0.Instrument < audio.InstrumentPhaser {
			sfxNote0.Instrument += 1
			playSfx0()
		}
		if pi.Btnp(pi.Left) && sfxNote0.Instrument > 0 {
			sfxNote0.Instrument -= 1
			playSfx0()
		}
	}

	ebitengine.MustRun()
}

func playSfx0() {
	audio.Sync() // first send changed sfx to audio system
	audio.Play(0, 0, 0, 0)
}
