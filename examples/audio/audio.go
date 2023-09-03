package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio"
	"github.com/elgopher/pi/ebitengine"
)

func main() {
	sfx := audio.SoundEffect{
		Notes: [32]audio.Note{
			{
				Pitch:      audio.PitchG3,
				Instrument: audio.InstrumentNoise,
				Volume:     7,
			},
		},
		Speed: 255,
	}
	audio.SetSfx(0, sfx)

	pi.Update = func() {
		if pi.Btnp(pi.X) {
			audio.Sfx(0, audio.Channel0, 0, 0)
		}

		if pi.Btnp(pi.Up) && sfx.Notes[0].Pitch < 255 {
			sfx.Notes[0].Pitch += 1
			audio.SetSfx(0, sfx)
			audio.Sfx(0, audio.Channel0, 0, 0)
		}
		if pi.Btnp(pi.Down) && sfx.Notes[0].Pitch > 0 {
			sfx.Notes[0].Pitch -= 1
			audio.SetSfx(0, sfx)
			audio.Sfx(0, audio.Channel0, 0, 0)
		}
		if pi.Btnp(pi.Right) && sfx.Notes[0].Instrument < audio.InstrumentPhaser {
			sfx.Notes[0].Instrument += 1
			audio.SetSfx(0, sfx)
			audio.Sfx(0, audio.Channel0, 0, 0)
		}
		if pi.Btnp(pi.Left) && sfx.Notes[0].Instrument > 0 {
			sfx.Notes[0].Instrument -= 1
			audio.SetSfx(0, sfx)
			audio.Sfx(0, audio.Channel0, 0, 0)
		}
	}

	ebitengine.MustRun()
}
