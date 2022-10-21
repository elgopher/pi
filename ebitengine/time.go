// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"time"

	"github.com/elgopher/pi/vm"
)

var lastTime time.Time

func updateTime() {
	now := time.Now()
	timePassed := now.Sub(lastTime)
	lastTime = now
	vm.TimeSeconds += float64(timePassed) / float64(time.Second)
}
