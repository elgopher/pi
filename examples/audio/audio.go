package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio"
	"github.com/elgopher/pi/ebitengine"
)

func main() {
	audio.SetSfx(0, audio.SoundEffect{
		Notes: [32]audio.Note{
			{
				Pitch:      audio.PitchC4,
				Instrument: audio.InstrumentOrgan,
				Volume:     7,
			},
		},
		Speed: 1,
	})

	pi.Update = func() {
		if pi.Btnp(pi.X) {
			audio.Sfx(0, 0, 0, 31)
		}
	}

	ebitengine.MustRun()
}
