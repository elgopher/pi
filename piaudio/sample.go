// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

// NewSample creates a new *Sample from raw 8-bit mono PCM data.
// baseFreq specifies the base frequency of the sample in Hz.
func NewSample(data []int8, baseFreq Freq) *Sample {
	return &Sample{data: data, baseFreq: baseFreq}
}

// Sample represents raw 8-bit mono PCM audio data along with its base frequency.
type Sample struct {
	data     []int8
	baseFreq Freq
}

func (s *Sample) Data() []int8 {
	return s.data
}

func (s *Sample) Len() int {
	return len(s.data)
}

func (s *Sample) BaseFreq() float64 {
	return s.baseFreq
}
