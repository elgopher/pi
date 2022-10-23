package inspector

import (
	"fmt"
	"math"

	"github.com/elgopher/pi"
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

func Update() {
	x, y := pi.MousePos()

	if pi.MouseBtn(pi.MouseLeft) && !distance.measuring {
		distance.measuring = true
		distance.startX, distance.startY = x, y
		fmt.Printf("Measuring started at (%d, %d)\n", x, y)
	} else if !pi.MouseBtn(pi.MouseLeft) && distance.measuring {
		distance.measuring = false
		dist, width, height := calcDistance()
		fmt.Printf("Measuring stopped at (%d, %d). Distance is: %f, width: %d, height: %d.\n", x, y, dist, width, height)
	}
}
