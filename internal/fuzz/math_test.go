// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func FuzzMinInt(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y int) {
		pi.MinInt(x, y)
	})
}

func FuzzMaxInt(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y int) {
		pi.MaxInt(x, y)
	})
}

func FuzzMidInt(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y, z int) {
		pi.MidInt(x, y, z)
	})
}

func FuzzMid(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y, z float64) {
		_ = pi.Mid(x, y, z)
	})
}
