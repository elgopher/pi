// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package snapshot

import (
	"github.com/elgopher/pi"
)

var snapshot []byte

func Take() {
	screen := pi.Scr()
	if snapshot == nil {
		snapshot = make([]byte, len(screen.Pix))
	}
	copy(snapshot, screen.Pix)
}

func Draw() {
	copy(pi.Scr().Pix, snapshot)
}
