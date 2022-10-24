// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/mem"
)

var (
	gamePaused     bool
	timeWhenPaused float64
)

func pauseGame() {
	gamePaused = true
	timeWhenPaused = mem.TimeSeconds
	snapshot.Take()
}

func resumeGame() {
	gamePaused = false
	mem.TimeSeconds = timeWhenPaused
	snapshot.Draw()
}
