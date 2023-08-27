// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import "math"

const (
	SampleRate = 22050
	Channels   = 4
)

type Samples [Channels]float64

func (s Samples) Sum() float64 {
	sum := 0.0
	for _, channelSample := range s {
		sum += channelSample
	}

	return sum
}

// PitchToFreq returns frequency in Hz. See audio.Pitch.
func PitchToFreq(pitch int) float64 {
	const pitchA2 = 33
	multiplier := math.Pow(2, float64(pitch-pitchA2)/12)
	return 440 * multiplier
}
