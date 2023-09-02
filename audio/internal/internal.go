// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

const (
	SampleRate = 22050
	Channels   = 4
)

type SampleChannels [Channels]float64

func (s SampleChannels) Sum() float64 {
	sum := 0.0
	for _, channelSample := range s {
		sum += channelSample
	}

	return sum
}
