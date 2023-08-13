// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
	"github.com/elgopher/pi/devtools/internal/snapshot"
)

type Shape struct {
	start   pi.Position
	started bool
	draw    func(x0, y0, x1, y1 int, color byte) string
	icon    byte
}

func (l *Shape) Update() {
	switch {
	case pi.MouseBtnp(pi.MouseLeft) && !l.started:
		l.start = pi.MousePos
		l.started = true
	case !pi.MouseBtn(pi.MouseLeft) && l.started:
		l.started = false
		snapshot.Draw()
		command := l.draw(l.start.X, l.start.Y, pi.MousePos.X, pi.MousePos.Y, FgColor)
		fmt.Println(command)
		snapshot.Take()
	}
}

func (l *Shape) Draw() {
	x, y := pi.MousePos.X, pi.MousePos.Y
	icons.Draw(x, y, FgColor, icons.Pointer)
	icons.Draw(x+2, y+2, FgColor, l.icon)
	if pi.MouseBtn(pi.MouseLeft) && l.started {
		l.draw(l.start.X, l.start.Y, x, y, FgColor)
	}
}

func (l *Shape) Icon() byte {
	return l.icon
}
