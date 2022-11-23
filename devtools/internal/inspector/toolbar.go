// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
	"github.com/elgopher/pi/key"
)

const toolbarWidth = 28
const toolbarHeight = 4

var toolbar Toolbar

type Toolbar struct {
	visible         bool
	pos             pi.Position
	toolHighlighted byte
}

func (t *Toolbar) toggle() {
	t.visible = !t.visible
	t.pos = pi.MousePosition
	t.pos.X -= toolbarWidth/2 - 2
	if t.pos.X < 0 {
		t.pos.X = 0
	}
	scr := pi.Scr()
	if t.pos.X+toolbarWidth > scr.Width() {
		t.pos.X = scr.Width() - toolbarWidth - 1
	}
	t.pos.Y -= toolbarHeight + 2
	if t.pos.Y < 0 {
		t.pos.Y = pi.MousePosition.Y + 6
	}
}

func (t *Toolbar) hide() {
	t.visible = false
	t.pos = pi.Position{}
	t.toolHighlighted = 0
}

func (t *Toolbar) update() {
	if t.visible {
		x, y := pi.MousePos()
		sx, sy := x-t.pos.X, y-t.pos.Y
		mouseOverToolbar := sx >= 0 && sy >= 0 && sx <= toolbarWidth && sy <= toolbarHeight
		if mouseOverToolbar {
			t.toolHighlighted = byte((sx-1)/4) + 1
		} else {
			t.toolHighlighted = 0
		}

		switch {
		case pi.MouseBtnp(pi.MouseLeft) && mouseOverToolbar:
			selectTool(t.toolHighlighted)
			t.hide()
		case pi.MouseBtnp(pi.MouseRight) || key.Btn(key.Esc):
			t.hide()
		}
	} else {
		if pi.MouseBtnp(pi.MouseRight) {
			t.toggle()
		}
	}
}

func (t *Toolbar) draw() {
	if t.visible {
		x := t.pos.X
		y := t.pos.Y
		pi.RectFill(x, y, x+toolbarWidth, y+toolbarHeight, BgColor)
		icons.Draw(x+1, y+1, FgColor,
			icons.MeasureTool, icons.PsetTool, icons.LineTool, icons.RectTool, icons.RectFillTool, icons.CircTool, icons.CircFillTool)

		if t.toolHighlighted > 0 {
			toolX := x + int((t.toolHighlighted-1)*4)
			pi.RectFill(toolX, y, toolX+4, y+4, FgColor)
			icons.Draw(toolX+1, y+1, BgColor, t.toolHighlighted)
		}

		t.drawPointer()
	}
}

func (t *Toolbar) drawPointer() {
	if t.toolHighlighted == 0 {
		x, y := pi.MousePos()
		icons.Draw(x, y, FgColor, icons.Pointer)
	}
}