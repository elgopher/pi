// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

// NewSample creates a new *Sample from raw 8-bit mono PCM data.
func NewSample(data []int8) *Sample {
	return &Sample{data: data}
}

// Sample represents raw 8-bit mono PCM audio data along with its base frequency.
type Sample struct {
	data []int8
}

func (s *Sample) Data() []int8 {
	return s.data
}

func (s *Sample) Len() int {
	return len(s.data)
}
