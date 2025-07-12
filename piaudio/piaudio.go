// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piaudio provides a simple tracker-style audio playback API
// for playing 8-bit mono PCM samples on 4 Amiga-style channels
// with tick-based scheduling.
//
// The design mimics classic MOD/XM-style playback where audio commands
// are scheduled in discrete ticks (TPS = ticks per second).
// This allows precise control of sample playback, pitch, volume,
// and timing effects with per-tick granularity.
package piaudio

// NewSample creates a new Sample from raw 8-bit mono PCM data.
func NewSample(data []byte) *Sample {
	return &Sample{data: data}
}

// Sample represents an 8-bit mono PCM sample.
type Sample struct {
	data []byte
	sent bool
}

// Send uploads the sample data to the audio backend.
func (s *Sample) Send() {
	s.sent = true
	Backend.Send(s)
}

// Freq represents frequency in Hz.
type Freq = float64

// Predefined common musical pitches in Hz.
const (
	A4 Freq = 440.0  // Standard tuning reference
	C4 Freq = 261.63 // Middle C
	E4 Freq = 329.63
	G4 Freq = 392.00
	A5 Freq = 880.0
)

// Loop defines sample loop parameters.
type Loop struct {
	Start    int // Start position in samples
	End      int // End position in samples
	LoopMode LoopMode
}

// LoopMode describes the type of sample looping behavior.
type LoopMode string

const (
	LoopModeNone     LoopMode = ""          // No looping
	LoopModeForward  LoopMode = "forward"   // Simple forward looping
	LoopModePingPong LoopMode = "ping-pong" // Alternating forward and backward loop
	LoopModeOneShot  LoopMode = "one-shot"  // Play once without looping
)

// TPS defines the number of audio ticks per second.
//
// This value defines the resolution of the audio scheduler,
// and typically matches tracker tempo (e.g. 50 for classic Amiga-style).
var TPS = 50

// Tick represents a unit of time for scheduling audio events.
//
// All audio events are planned in terms of ticks,
// allowing precise timing for playback and parameter changes.
type Tick = int

// CurrentTick is the current global tick counter.
//
// The CurrentTick is set by backend.
var CurrentTick Tick = 0

// Play schedules playback of a Sound on the specified channel
// with delay in ticks.
//
// The delay is relative to the CurrentTick:
// delay = 0 means play immediately, delay = 10 means play after 10 ticks.
//
// Example:
//
//	piaudio.Play(piaudio.Chan0, sound, 0)    // play now
//	piaudio.Play(piaudio.Chan1, sound, 5)    // play in 5 ticks
func Play(ch Chan, sound Sound, delay Tick) {
	if !sound.Sample.sent {
		sound.Sample.Send() // send only the first time
		sound.Sample.sent = true
	}

	Backend.Play(ch, sound, delay)
}

// Chan represents an audio channel number.
type Chan = uint8

const (
	Chan0 Chan = 0 // Left speaker
	Chan1 Chan = 1 // Left speaker
	Chan2 Chan = 2 // Right speaker
	Chan3 Chan = 3 // Right speaker
)

// Stop schedules stopping playback on the specified channel
// with an optional delay in ticks.
//
// delay = 0 means stop immediately, delay = 5 means stop after 5 ticks.
func Stop(ch Chan, delay Tick) {
	Backend.Stop(ch, delay)
}

// SetPitch schedules a pitch change on the specified channel
// with an optional delay in ticks.
//
// Pitch is specified in Hz.
// delay = 0 means change immediately, delay = 3 means after 3 ticks
func SetPitch(ch Chan, pitch Freq, delay Tick) {
	Backend.SetPitch(ch, pitch, delay)
}

// SetVolume schedules a volume change on the specified channel
// with an optional delay in ticks.
//
// Volume is specified as a float in the range 0.0–1.0.
// delay = 0 means change immediately, delay = 2 means after 2 ticks.
func SetVolume(ch Chan, vol float64, delay Tick) {
	Backend.SetVolume(ch, vol, delay)
}
