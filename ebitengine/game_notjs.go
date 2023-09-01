// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package ebitengine

import "github.com/hajimehoshi/ebiten/v2"

// on non-js operating systems, initialization takes only few milliseconds and there is no user action needed.
// Therefore, there is no need to draw anything.
func (e *game) drawNotReady(*ebiten.Image) {}
