// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

const color = 7

func FuzzPset(f *testing.F) {
	pi.Reset()
	pi.SetScreenSize(16, 16)
	f.Fuzz(func(t *testing.T, x, y int) {
		pi.Pset(x, y, color)
	})
}

func FuzzPget(f *testing.F) {
	pi.Reset()
	pi.SetScreenSize(16, 16)
	f.Fuzz(func(t *testing.T, x, y int) {
		pi.Pget(x, y)
	})
}

func FuzzSprSizeFlip(f *testing.F) {
	pi.Reset()
	pi.SetScreenSize(16, 16)
	f.Fuzz(func(t *testing.T, n, x, y int, w, h float64, flipX, flipY bool) {
		pi.SprSizeFlip(n, x, y, w, h, flipX, flipY)
	})
}
