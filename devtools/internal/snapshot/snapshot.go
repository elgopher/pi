// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package snapshot

import (
	"github.com/elgopher/pi/mem"
)

var snapshot []byte

func Take() {
	if snapshot == nil {
		snapshot = make([]byte, len(mem.ScreenData))
	}
	copy(snapshot, mem.ScreenData)
}

func Draw() {
	copy(mem.ScreenData, snapshot)
}
