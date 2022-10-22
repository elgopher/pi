// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/vm"
)

var (
	paused         bool
	timeWhenPaused float64
)

func pauseGame() {
	paused = true
	timeWhenPaused = vm.TimeSeconds
	snapshot.Take()
}

func resumeGame() {
	paused = false
	vm.TimeSeconds = timeWhenPaused
}
