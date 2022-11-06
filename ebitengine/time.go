// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"time"

	"github.com/elgopher/pi"
)

var lastTime time.Time

func updateTime() {
	if lastTime.IsZero() {
		lastTime = time.Now()
		pi.TimeSeconds = 0
		return
	}

	now := time.Now()
	timePassed := now.Sub(lastTime)
	lastTime = now
	pi.TimeSeconds += float64(timePassed) / float64(time.Second)
}
