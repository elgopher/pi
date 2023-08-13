// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
	"github.com/elgopher/pi/devtools/internal/snapshot"
)

type Set struct {
	running bool
}

func (p *Set) Update() {
	x, y := pi.MousePos.X, pi.MousePos.Y
	leftp := pi.MouseBtnp(pi.MouseLeft) && !p.running
	left := pi.MouseBtn(pi.MouseLeft) && p.running
	if (leftp || left) && pixelColorAtMouseCoords != FgColor {
		p.running = true
		snapshot.Draw()
		pi.Set(x, y, FgColor)
		fmt.Printf("pi.Set(%d, %d, %d)\n", x, y, FgColor)
		snapshot.Take()
	}
}

func (p *Set) Draw() {
	pi.Set(pi.MousePos.X, pi.MousePos.Y, FgColor)
}

func (p *Set) Icon() byte {
	return icons.SetTool
}
