// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package snapshot

import (
	"github.com/elgopher/pi"
)

var snapshot []byte
var history [][]byte

const historyLimit = 128

func Take() {
	if snapshot != nil {
		s := make([]byte, len(snapshot))
		copy(s, snapshot)
		history = append(history, s)
		if len(history) > historyLimit {
			history = history[1:]
		}
	}

	screen := pi.Scr()
	pix := screen.Pix()
	if snapshot == nil {
		snapshot = make([]byte, len(pix))
	}
	copy(snapshot, pix)
}

func Draw() {
	copy(pi.Scr().Pix(), snapshot)
}

func Undo() {
	if len(history) == 0 {
		return
	}

	last := history[len(history)-1]
	history = history[:len(history)-1]
	snapshot = last
}
