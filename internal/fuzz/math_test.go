// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func FuzzMidInt(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y, z int) {
		pi.Mid(x, y, z)
	})
}

func FuzzMid(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y, z float64) {
		_ = pi.Mid(x, y, z)
	})
}
