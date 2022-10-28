// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"fmt"
	"math"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/key"
)

var distance struct {
	measuring      bool
	startX, startY int
}

func calcDistance() (dist float64, width, height int) {
	x, y := pi.MousePos()

	width = x - distance.startX
	if width < 0 {
		width *= -1
	}

	height = y - distance.startY
	if height < 0 {
		height *= -1
	}
	dist = math.Sqrt(float64(width*width + height*height))
	return
}

var helpShown bool

func Update() {
	if !helpShown {
		helpShown = true
		fmt.Println("Press right mouse button to show toolbar.")
		fmt.Println("Press P to take screenshot.")
	}

	if !toolbar.visible {
		tool.Update()
	}
	toolbar.update()

	if key.Btn(key.Ctrl) && key.Btnp(key.Z) {
		snapshot.Undo()
	}
}
