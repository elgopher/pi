// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"fmt"
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
	case icons.SetTool:
		tool = &Set{}
	case icons.LineTool:
		tool = &Shape{
			draw: drawShape("Line", pi.Line),
			icon: icons.LineTool,
		}
	case icons.RectTool:
		tool = &Shape{
			draw: drawShape("Rect", pi.Rect),
			icon: icons.RectTool,
		}
	case icons.RectFillTool:
		tool = &Shape{
			draw: drawShape("RectFill", pi.RectFill),
			icon: icons.RectFillTool,
		}
	case icons.CircTool:
		tool = &Shape{
			draw: drawCirc("Circ", pi.Circ),
			icon: icons.CircTool,
		}
	case icons.CircFillTool:
		tool = &Shape{
			draw: drawCirc("CircFill", pi.CircFill),
			icon: icons.CircFillTool,
		}
	}
}

func drawShape(name string, f func(x0, y0, x1, y1 int, color byte)) func(x0, y0, x1, y1 int, color byte) string {
	return func(x0, y0, x1, y1 int, color byte) string {
		f(x0, y0, x1, y1, color)
		command := fmt.Sprintf("pi.%s(%d, %d, %d, %d, %d)", name, x0, y0, x1, y1, color)
		return command
	}
}

func drawCirc(name string, f func(cx, cy, r int, color byte)) func(x0 int, y0 int, x1 int, y1 int, color byte) string {
	return func(x0, y0, x1, y1 int, color byte) string {
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

		command := fmt.Sprintf("pi.%s(%d, %d, %d, %d)", name, cx, cy, r, color)
		return command
	}
}
