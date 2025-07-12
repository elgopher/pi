// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

// Sound describes all playback parameters for scheduling a sample on a channel.
type Sound struct {
	// Sample is the PCM data to play.
	Sample *Sample
	// BaseFreq is the base frequency of the sample in Hz.
	BaseFreq Freq
	// Loop defines the loop region and mode.
	Loop Loop
	// Pitch is the playback frequency in Hz.
	Pitch Freq
	// Offset is the start position in the sample data (in samples).
	// Used for sample offset effects (e.g. 9xx in tracker).
	Offset int
	// Volume is the playback volume (0.0–1.0).
	Volume float64
}

// WithSample returns a copy of the Sound with the specified Sample.
func (s Sound) WithSample(sample *Sample) Sound {
	s.Sample = sample
	return s
}

// WithBaseFreq returns a copy of the Sound with the specified base frequency.
func (s Sound) WithBaseFreq(baseFreq Freq) Sound {
	s.BaseFreq = baseFreq
	return s
}

// WithLoop returns a copy of the Sound with the specified Loop settings.
func (s Sound) WithLoop(loop Loop) Sound {
	s.Loop = loop
	return s
}

// WithPitch returns a copy of the Sound with the specified playback Pitch.
func (s Sound) WithPitch(pitch Freq) Sound {
	s.Pitch = pitch
	return s
}

// WithOffset returns a copy of the Sound with the specified sample Offset.
func (s Sound) WithOffset(offset int) Sound {
	s.Offset = offset
	return s
}

// WithVolume returns a copy of the Sound with the specified playback Volume.
func (s Sound) WithVolume(volume float64) Sound {
	s.Volume = volume
	return s
}
