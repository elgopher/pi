// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/snapshot"
)

var (
	gamePaused     bool
	timeWhenPaused float64
)

func pauseGame() {
	gamePaused = true
	timeWhenPaused = pi.TimeSeconds
	snapshot.Take()
}

func resumeGame() {
	gamePaused = false
	pi.TimeSeconds = timeWhenPaused
	snapshot.Draw()
}
