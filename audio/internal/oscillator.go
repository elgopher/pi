// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import "math"

type Oscillator struct {
	Func   func(pos float64) float64
	FreqHz float64
	Pos    float64
}

func (o *Oscillator) NextSample() (v float64) {
	v = o.Func(o.Pos)
	o.Pos += o.FreqHz / SampleRate

	return
}

const volume = 0.7

func Triangle(pos float64) float64 {
	phase := math.Mod(pos, 1)
	value := math.Abs(phase*2-1)*2 - 1

	return value * volume
}

// Organ is triangle / 2
func Organ(pos float64) float64 {
	pos = pos * 4

	return (math.Abs(math.Mod(pos, 2)-1) - 0.5 +
		(math.Abs((math.Mod(pos*0.5, 2))-1)-0.5)/2 - 0.1) * volume
}
