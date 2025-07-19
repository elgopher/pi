// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piaudio provides a low-level API for sound generation,
// inspired by the Paula audio chip used in Amiga computers.
//
// In piaudio, sound is generated in 4 independent channels,
// which are then mixed to stereo output:
// channels 0 and 3 are mixed to the left speaker,
// and channels 1 and 2 to the right.
//
// Each channel plays a single audio sample at a time,
// which means you can play 4 different samples simultaneously.
// A channel, in addition to holding the current sample,
// also stores parameters such as pitch, volume, and loop,
// which affect how the sound is generated.
// These parameters can be updated in real-time
// to play sound effects and music.
//
// The piaudio package supports scheduling operations in advance.
// It allows setting the sample to play or changing pitch and volume
// with a specified delay parameter.
// Programmers should plan commands tens of milliseconds ahead
// to avoid audio glitches caused by CPU blocking, e.g. by the garbage collector.
//
// piaudio only supports samples encoded in 8-bit mono PCM format.
// This means each audio sample is stored as a single signed integer (int8)
// in the range -128 to 127, representing the audio signal amplitude at a point in time.
// "Mono" means the sample contains just one channel (no left/right separation),
// so all sound data refers to a single audio source.
// PCM (Pulse-Code Modulation) is an uncompressed format
// where amplitude values are stored directly,
// enabling very simple and fast playback and processing.
//
// This format is easy to generate and modify in code
// and is well suited to retro-style constraints,
// mimicking the behavior of vintage audio hardware.
//
// To import your own audio samples into piaudio,
// decode them into a *Sample object using the DecodeWav or DecodeRaw functions,
// or provide raw data directly to NewSample.
// Once created, a *Sample can be loaded into the backend using LoadSample.
package piaudio

// Chan represents an audio channel number.
type Chan = uint8

const (
	Chan1 Chan = 1 // Left speaker
	Chan2 Chan = 2 // Right speaker
	Chan3 Chan = 4 // Right speaker
	Chan4 Chan = 8 // Left speaker
)

// Time set by the audio backend for synchronization.
var Time float64

// LoadSample loads a sample into the backend. Call this before playing the sample.
func LoadSample(sample *Sample) { Backend.LoadSample(sample) }

// UnloadSample removes a sample from the backend. Call this to free memory when it's no longer needed.
func UnloadSample(sample *Sample) { Backend.UnloadSample(sample) }

// Play plays a given sample on selected channel(s).
//
// A pitch of 1.0 plays the sample at its original speed.
// Values below 1.0 slow down the sample and lower its pitch (e.g., 0.5 = one octave lower),
// while values above 1.0 speed it up and raise the pitch.
//
// Vol is 0-1. 1 is the loudest.
func Play(ch Chan, sample *Sample, pitch, vol float64) {
	delay := 0.0
	ClearChan(ch, delay)
	SetSample(ch, sample, 0, delay)
	SetLoop(ch, 0, sample.Len(), LoopNone, delay)
	SetPitch(ch, pitch, delay)
	SetVolume(ch, vol, delay)
}

// SetSample schedules playback of the sample to take effect after the specified delay.
//
// Initial sample is nil, offset is 0.
func SetSample(ch Chan, sample *Sample, offset int, delay float64) {
	Backend.SetSample(ch, sample, offset, delay)
}

// SetLoop schedules the loop configuration to take effect after the specified delay.
//
// This function can also be used to limit the sample's length. Just call
// SetLoop with loopType = LoopNone and provide the desired length.
//
// Initial start is 0, length is 2_147_483_647 and loopType is LoopNone
func SetLoop(ch Chan, start int, length int, loopType LoopType, delay float64) {
	Backend.SetLoop(ch, start, length, loopType, delay)
}

// SetPitch schedules the pitch change to take effect after the specified delay.
// A pitch of 1.0 plays the sample at its original speed.
// Values below 1.0 slow down the sample and lower its pitch (e.g., 0.5 = one octave lower),
// while values above 1.0 speed it up and raise the pitch.
//
// Initial pitch is 1.0
func SetPitch(ch Chan, pitch float64, delay float64) {
	Backend.SetPitch(ch, pitch, delay)
}

// SetVolume schedules the volume change to take effect after the specified delay.
//
// Initial vol is 1.0
func SetVolume(ch Chan, vol float64, delay float64) {
	Backend.SetVolume(ch, vol, delay)
}

// ClearChan removes all scheduled operations for the channel after the specified delay.
func ClearChan(ch Chan, delay float64) {
	Backend.ClearChan(ch, delay)
}
