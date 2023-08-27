package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio"
	"github.com/elgopher/pi/ebitengine"
)

var playing bool

func main() {
	audio.SetSfx(0, audio.SoundEffect{
		Notes: [32]audio.Note{
			{
				Pitch:      audio.PitchC4,
				Instrument: audio.InstrumentOrgan,
				Volume:     7,
			},
		},
		Speed: 10,
	})

	pi.Update = func() {
		if pi.Btnp(pi.X) {
			if !playing {
				audio.Sfx(0, audio.Channel0, 0, 31)
				playing = true
			} else {
				audio.Sfx(-1, audio.Channel0, 0, 0)
				playing = false
			}
		}
	}

	ebitengine.MustRun()
}
