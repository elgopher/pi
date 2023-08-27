package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio"
	"github.com/elgopher/pi/ebitengine"
)

var playing bool

func main() {
	sfx := audio.SoundEffect{
		Notes: [32]audio.Note{
			{
				Pitch:      audio.PitchC2,
				Instrument: audio.InstrumentOrgan,
				Volume:     7,
			},
		},
		Speed: 10,
	}
	audio.SetSfx(0, sfx)

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

		if playing {
			if pi.Btnp(pi.Up) && sfx.Notes[0].Pitch < 255 {
				sfx.Notes[0].Pitch += 1
				audio.SetSfx(0, sfx)
				audio.Sfx(0, audio.Channel0, 0, 31)
			}
			if pi.Btnp(pi.Down) && sfx.Notes[0].Pitch > 0 {
				sfx.Notes[0].Pitch -= 1
				audio.SetSfx(0, sfx)
				audio.Sfx(0, audio.Channel0, 0, 31)
			}
		}
	}

	ebitengine.MustRun()
}
