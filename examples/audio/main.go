package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

func main() {
	audio := pi.Audio()
	effect := &audio.Effects[0]

	note := &effect.Notes[0]
	note.Volume = 7
	note.Pitch = 444
	note.Wave = pi.WaveTriangle

	pi.Sfx(0)
	ebitengine.MustRun()
}
