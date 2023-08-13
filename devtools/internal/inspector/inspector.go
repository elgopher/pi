// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"github.com/elgopher/pi"
)

var isBarOnTop bool

func moveBarIfNeeded() {
	switch {
	case isBarOnTop && pi.MousePos.Y <= 12:
		isBarOnTop = false
	case !isBarOnTop && pi.MousePos.Y >= pi.Scr().Height()-12:
		isBarOnTop = true
	}
}
