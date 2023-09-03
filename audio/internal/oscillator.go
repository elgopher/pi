// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"math"
	"math/rand"
)

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

func Triangle(pos float64) float64 {
	// Based on https://github.com/picolove
	phase := math.Mod(pos, 1)
	value := math.Abs(phase*2-1)*2 - 1

	return value * 0.45
}

func TiltedSaw(pos float64) float64 {
	// Based on https://github.com/picolove
	phase := math.Mod(pos, 1)
	var v float64
	if phase < 0.875 {
		v = phase * 16 / 7
	} else {
		v = (1 - phase) * 16
	}
	return (v - 1) * 0.45
}

func Saw(pos float64) float64 {
	// Based on https://github.com/picolove
	phase := math.Mod(pos, 1)
	return (phase - 0.5) * 0.65
}

func Square(pos float64) float64 {
	// Based on https://github.com/picolove
	phase := math.Mod(pos, 1)
	v := -1.0
	if phase < 0.5 {
		v = 1.0
	}
	return v / 4.0
}

func Pulse(pos float64) float64 {
	// Based on https://github.com/picolove
	phase := math.Mod(pos, 1)
	v := -1.0
	if phase < 0.3125 {
		v = 1.0
	}
	return v / 4.0
}

// Organ is triangle / 2
func Organ(pos float64) float64 {
	// Based on https://github.com/picolove
	pos = pos * 4

	v := math.Abs(math.Mod(pos, 2)-1) - 0.5 +
		(math.Abs((math.Mod(pos*0.5, 2))-1)-0.5)/2 - 0.1

	return v * 0.55
}

func Noise() func(pos float64) float64 {
	// Based on https://github.com/picolove
	lastPos := 0.0
	sample := 0.0
	lsample := 0.0
	const tscale = 2489.0158697766 / SampleRate // 2489.0158697766 is freq for D#5

	return func(pos float64) float64 {
		scale := (pos - lastPos) / tscale
		lsample = sample
		sample = (lsample + scale*(rand.Float64()*2-1)) / (1 + scale)
		lastPos = pos
		v := math.Min(math.Max((lsample+sample)*4/3*(1.75-scale), -1), 1)
		return v * 0.5
	}
}

func Phaser(pos float64) float64 {
	// Based on https://github.com/picolove
	pos = pos * 2

	v := math.Abs(math.Mod(pos, 2)-1) - 0.5 +
		(math.Abs(math.Mod(pos*127/128, 2)-1)-0.5)/2 - 0.25

	return v * 0.7
}
