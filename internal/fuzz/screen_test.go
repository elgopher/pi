// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

const color = 7

func FuzzPixMap_Set(f *testing.F) {
	pixMap := pi.NewPixMap(16, 16)
	f.Fuzz(func(t *testing.T, x, y int) {
		pixMap.Set(x, y, color)
	})
}

func FuzzPixMap_Get(f *testing.F) {
	pixMap := pi.NewPixMap(16, 16)
	f.Fuzz(func(t *testing.T, x, y int) {
		pixMap.Get(x, y)
	})
}

func FuzzSprSizeFlip(f *testing.F) {
	pi.Reset()
	pi.SetScreenSize(16, 16)
	f.Fuzz(func(t *testing.T, n, x, y int, w, h float64, flipX, flipY bool) {
		pi.SprSizeFlip(n, x, y, w, h, flipX, flipY)
	})
}
