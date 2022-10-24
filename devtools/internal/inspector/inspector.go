// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/mem"
)

var isBarOnTop bool

func moveBarIfNeeded() {
	_, mouseY := pi.MousePos()
	switch {
	case isBarOnTop && mouseY <= 12:
		isBarOnTop = false
	case !isBarOnTop && mouseY >= mem.ScreenHeight-12:
		isBarOnTop = true
	}
}
