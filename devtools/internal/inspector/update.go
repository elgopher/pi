// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
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
	width = pi.MousePos.X - distance.startX
	if width < 0 {
		width *= -1
	}

	height = pi.MousePos.Y - distance.startY
	if height < 0 {
		height *= -1
	}
	dist = math.Sqrt(float64(width*width + height*height))
	return
}

func Update() {

	if !toolbar.visible {
		tool.Update()
	}
	toolbar.update()

	if key.Btn(key.Ctrl) && key.Btnp(key.Z) {
		snapshot.Undo()
	}
}
