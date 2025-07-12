// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import "github.com/elgopher/pi/piaudio"

type audioBackend struct{}

func (audioBackend) Send(*piaudio.Sample) {
}

func (audioBackend) Play(piaudio.Chan, piaudio.Sound, piaudio.Tick) {
}

func (audioBackend) Stop(piaudio.Chan, piaudio.Tick) {
}

func (audioBackend) SetPitch(piaudio.Chan, piaudio.Freq, piaudio.Tick) {
}

func (audioBackend) SetVolume(piaudio.Chan, float64, piaudio.Tick) {
}
