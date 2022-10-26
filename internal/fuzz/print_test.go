// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func FuzzPrint(f *testing.F) {
	pi.Reset()
	pi.SetScreenSize(16, 24)
	f.Fuzz(func(t *testing.T, x, y int) {
		pi.Print("A", x, y, color)
	})
}
