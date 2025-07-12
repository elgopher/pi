// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

var Backend interface {
	Send(*Sample)
	Play(_ Chan, sound Sound, delay Tick)
	Stop(_ Chan, delay Tick)
	SetPitch(_ Chan, pitch Freq, delay Tick)
	SetVolume(_ Chan, vol float64, delay Tick)
} = panicBackend{}

type panicBackend struct{}

func (panicBackend) Send(*Sample) {
	panic("Cannot send sample: backend not set. Please send the sample after running the game")
}

func (panicBackend) Play(Chan, Sound, Tick) {
	panic("Cannot play sound: backend not set. Please play the sound after running the game")
}

func (panicBackend) Stop(Chan, Tick) {
	panic("Cannot stop audio channel: backend not set. Please stop the sound after running the game")
}

func (panicBackend) SetPitch(Chan, Freq, Tick) {
	panic("Cannot set channel pitch: backend not set. Please update pitch after running the game")
}

func (panicBackend) SetVolume(Chan, float64, Tick) {
	panic("Cannot set volume pitch: backend not set. Please update volume after running the game")
}
