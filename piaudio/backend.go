// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

// Backend variable is set by the backend (e.g. piebiten) after starting
// the game, for example with piebiten.Run.
var Backend BackendInterface = panicBackend{}

// BackendInterface is implemented by backends, e.g. piebiten.
type BackendInterface interface {
	// LoadSample loads a sample into the backend. Call this before playing the sample.
	LoadSample(sample *Sample)

	// UnloadSample removes a sample from the backend. Call this to free memory when it's no longer needed.
	UnloadSample(sample *Sample)

	// SetSample schedules playback of the sample to take effect after the specified delay.
	SetSample(ch Chan, sample *Sample, offset int, delay float64)

	// SetLoop schedules the loop configuration to take effect after the specified delay.
	SetLoop(_ Chan, start int, stop int, loop LoopType, delay float64)

	// SetPitch schedules the pitch change to take effect after the specified delay.
	// A pitch of 1.0 plays the sample at its original speed.
	// Values below 1.0 slow down the sample and lower its pitch (e.g., 0.5 = one octave lower),
	// while values above 1.0 speed it up and raise the pitch.
	SetPitch(_ Chan, pitch float64, delay float64)

	// SetVolume schedules the volume change to take effect after the specified delay.
	SetVolume(_ Chan, vol float64, delay float64)

	// ClearChan removes all scheduled operations for the channel after the specified delay.
	ClearChan(ch Chan, delay float64)
}

type panicBackend struct{}

func (p panicBackend) LoadSample(sample *Sample) {
	panic("cannot load sample: backend not set. Please call LoadSample only after starting the game")
}

func (p panicBackend) UnloadSample(sample *Sample) {
	panic("cannot unload sample: backend not set. Please call UnloadSample only after starting the game")
}

func (p panicBackend) SetSample(ch Chan, sample *Sample, offset int, delay float64) {
	panic("cannot set sample: backend not set. Please call SetSample only after starting the game")
}

func (p panicBackend) SetLoop(ch Chan, start, stop int, loop LoopType, delay float64) {
	panic("cannot set loop: backend not set. Please call SetLoop only after starting the game")
}

func (p panicBackend) ClearChan(ch Chan, delay float64) {
	panic("cannot clear channel: backend not set. Please call ClearChan only after starting the game")
}

func (p panicBackend) SetPitch(_ Chan, pitch float64, delay float64) {
	panic("cannot set pitch: backend not set. Please call SetPitch only after starting the game")
}

func (p panicBackend) SetVolume(_ Chan, vol float64, delay float64) {
	panic("cannot set volume: backend not set. Please call SetVolume only after starting the game")
}
