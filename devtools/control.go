// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/vm"
)

var (
	gamePaused     bool
	timeWhenPaused float64
)

func pauseGame() {
	gamePaused = true
	timeWhenPaused = vm.TimeSeconds
	snapshot.Take()
}

func resumeGame() {
	gamePaused = false
	vm.TimeSeconds = timeWhenPaused
	snapshot.Draw()
}
