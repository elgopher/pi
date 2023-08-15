package inspector

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
	"github.com/elgopher/pi/devtools/internal/snapshot"
)

type Spr struct {
	icon byte
	mode sprMode
}

type sprMode interface {
	update(*Spr)
	draw()
}

func (l *Spr) Update() {
	l.mode.update(l)
}

func (l *Spr) Draw() {
	l.mode.draw()
}

func (l *Spr) Icon() byte {
	return l.icon
}

type sprSelectionMode struct {
	lastMousePos pi.Position
	camera       pi.Position
}

func (m *sprSelectionMode) update(spr *Spr) {
	if pi.MouseBtn(pi.MouseMiddle) {
		zero := pi.Position{}
		if m.lastMousePos == zero {
			m.lastMousePos = pi.MousePos
		}
		dx := m.lastMousePos.X - pi.MousePos.X
		dy := m.lastMousePos.Y - pi.MousePos.Y

		m.camera.X -= dx
		minCameraX := -(pi.SprSheet().Width() - pi.Scr().Width())
		m.camera.X = pi.MidInt(minCameraX, m.camera.X, 0)

		m.camera.Y -= dy
		minCameraY := -(pi.SprSheet().Height() - pi.Scr().Height())
		m.camera.Y = pi.MidInt(minCameraY, m.camera.Y, 0)

		m.lastMousePos = pi.MousePos
	}

	if !pi.MouseBtn(pi.MouseMiddle) {
		m.lastMousePos.Reset()
	}

	if pi.MouseBtnp(pi.MouseLeft) {
		spr.mode = &sprDrawingMode{selectedSprite: m.spriteNo(pi.MousePos)}
	}
}

func (m *sprSelectionMode) draw() {
	pi.SprSheet().Copy(0, 0, pi.SprSheet().Width(), pi.SprSheet().Height(), pi.Scr(), m.camera.X, m.camera.Y)

	mouseX, mouseY := pi.MousePos.X, pi.MousePos.Y

	icons.Draw(mouseX, mouseY, FgColor, icons.Pointer)

	screen := pi.Scr()
	var barY int
	if !isBarOnTop {
		barY = screen.Height() - 7
	}

	pi.RectFill(0, barY, screen.Width(), barY+6, BgColor)

	spriteNo := fmt.Sprintf("SRITE %d", m.spriteNo(pi.MousePos))
	pi.Print(spriteNo, 1, barY+1, FgColor)
}

func (m *sprSelectionMode) spriteNo(mousePos pi.Position) int {
	x := mousePos.X - m.camera.X
	y := mousePos.Y - m.camera.Y
	cellX := x / 8
	cellY := y / 8
	return cellY*(pi.SprSheet().Width()/8) + cellX
}

type sprDrawingMode struct {
	selectedSprite int
}

func (m *sprDrawingMode) update(*Spr) {
	if pi.MouseBtnp(pi.MouseLeft) {
		fmt.Printf("pi.Spr(%d, %d, %d)\n", m.selectedSprite, pi.MousePos.X, pi.MousePos.Y)
		snapshot.Take()
	}
}

func (m *sprDrawingMode) draw() {
	mouseX, mouseY := pi.MousePos.X, pi.MousePos.Y
	pi.Spr(m.selectedSprite, mouseX, mouseY)
}
