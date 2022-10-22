// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package snapshot

import (
	"github.com/elgopher/pi/vm"
)

var snapshot []byte

func Take() {
	if snapshot == nil {
		snapshot = make([]byte, len(vm.ScreenData))
	}
	copy(snapshot, vm.ScreenData)
}

func Draw() {
	copy(vm.ScreenData, snapshot)
}
