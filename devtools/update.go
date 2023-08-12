// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"github.com/elgopher/pi/devtools/internal/inspector"
	"github.com/elgopher/pi/key"
)

func updateDevTools() {
	if key.Btnp(key.F12) {
		if !gamePaused {
			pauseGame()
		} else {
			resumeGame()
		}
	}

	if gamePaused {
		inspector.Update()
	}

	evaluateNextCommandFromTerminal()
}
