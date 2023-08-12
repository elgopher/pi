// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"time"

	"github.com/elgopher/pi"
)

var lastTime time.Time

func updateTime() {
	now := time.Now()
	timePassed := now.Sub(lastTime)
	lastTime = now
	pi.Time += float64(timePassed) / float64(time.Second)
}
