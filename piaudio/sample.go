// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

// NewSample creates a new *Sample from raw 8-bit mono PCM data.
// sampleRate determines how many times per second a sample value is read (in Hz).
func NewSample(data []int8, sampleRate uint16) *Sample {
	return &Sample{data: data, sampleRate: sampleRate}
}

// Sample represents raw 8-bit mono PCM audio data along with its sample rate.
type Sample struct {
	data       []int8
	sampleRate uint16
}

func (s *Sample) Data() []int8 {
	return s.data
}

func (s *Sample) Len() int {
	return len(s.data)
}

func (s *Sample) SampleRate() uint16 {
	return s.sampleRate
}
