// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"math"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
)

var tool Tool = &Measure{}

type Tool interface {
	Draw()
	Update()
	Icon() byte
}

func selectTool(icon byte) {
	switch icon {
	case icons.MeasureTool:
		tool = &Measure{}
	case icons.PsetTool:
		tool = &Pset{}
	case icons.LineTool:
		tool = &Shape{
			draw:     pi.Line,
			function: "Line",
			icon:     icons.LineTool,
		}
	case icons.RectTool:
		tool = &Shape{
			draw:     pi.Rect,
			function: "Rect",
			icon:     icons.RectTool,
		}
	case icons.RectFillTool:
		tool = &Shape{
			draw:     pi.RectFill,
			function: "RectFill",
			icon:     icons.RectFillTool,
		}
	case icons.CircTool:
		tool = &Shape{
			draw:     drawCirc(pi.Circ),
			function: "Circ",
			icon:     icons.CircTool,
		}
	case icons.CircFillTool:
		tool = &Shape{
			draw:     drawCirc(pi.CircFill),
			function: "CircFill",
			icon:     icons.CircFillTool,
		}
	}
}

func drawCirc(f func(cx, cy, r int, color byte)) func(x0 int, y0 int, x1 int, y1 int, color byte) {
	return func(x0, y0, x1, y1 int, color byte) {
		dx := x1 - x0
		cx := x0 + dx
		if dx < 0 {
			dx *= -1
		}
		dy := y1 - y0
		cy := y0 + dy
		if dy < 0 {
			dy *= -1
		}

		r := int(math.Sqrt(float64(dx*dx + dy*dy)))

		f(cx, cy, r, FgColor)
	}
}
