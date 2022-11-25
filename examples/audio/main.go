package main

import (
	"fmt"
	"time"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

func main() {
	audio := pi.Audio()
	effect := audio.SoundEffect(0) // returns a copy of SoundEffect

	note := &effect.Notes[0] // modify. No need to synchronize access
	note.Volume = 7
	note.Pitch = 444
	note.Wave = pi.WaveTriangle

	audio.SetSoundEffect(effect) // concurrency-safe update of sound effect

	pi.Sfx(0) // play

	channel := audio.Channel(0) // return a copy what is currently playing
	fmt.Println(channel.Playing, channel.Sfx, channel.Note)

	// TODO Better to use index parameter?
	sfxSlice := pi.SfxSlice{No: 0, NoteStart: 0, NoteEnd: 31}
	audio.SetChannel(pi.ChannelState{No: 0, Sfx: sfxSlice, Note: 0}) // instead of pi.Sfx(0)

	time.Sleep(time.Second)

	pi.SfxStop(0)

	ebitengine.MustRun() // backend runs audio.Read(buf) in a loop
}
