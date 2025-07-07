// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func TestStretch(t *testing.T) {
	// temporary test
	dst := pi.NewCanvas(16, 16)
	pi.SetDrawTarget(dst)

	src := pi.NewCanvas(8, 8)
	src.Clear(7)

	spr := pi.CanvasSprite(src)

	pi.Stretch(spr, 0, 0, 8, 8)
	pi.Stretch(spr, -1, 0, 8, 8)
	pi.Stretch(spr, 0, -1, 8, 8)
	pi.Stretch(spr, 16, 0, 8, 8)
	pi.Stretch(spr, 0, 16, 8, 8)

	pi.Stretch(spr.WithFlipX(true), 0, 0, 8, 8)
	pi.Stretch(spr.WithFlipY(true), 0, 0, 8, 8)

	pi.Stretch(spr.WithSize(0, 0), 0, 0, 8, 8)
}
