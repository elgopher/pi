// Example showing how to test pressed buttons of game controllers.
package main

import (
	"embed"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

//go:embed sprite-sheet.png
var resources embed.FS

// colors
const (
	left           = 10
	up             = 9
	down           = 13
	right          = 11
	btnO           = 4
	btnX           = 3
	dpadInactive   = 5
	buttonInactive = 8
	active         = 14

	yellow = 10
)

func main() {
	pi.Resources = resources
	pi.Draw = func() {
		pi.Cls()
		drawPlayerController(0, 2, 20)
		drawPlayerController(1, 2, 70)
		drawPlayerController(2, 67, 20)
		drawPlayerController(3, 67, 70)
	}
	pi.MustRun(ebitengine.Backend)
}

func drawPlayerController(player, x, y int) {
	pi.PalReset()
	// make all buttons inactive:
	pi.Pal(left, dpadInactive)
	pi.Pal(up, dpadInactive)
	pi.Pal(down, dpadInactive)
	pi.Pal(right, dpadInactive)
	pi.Pal(btnO, buttonInactive)
	pi.Pal(btnX, buttonInactive)

	// make pressed buttons active:
	if pi.BtnPlayer(pi.O, player) {
		pi.Pal(btnO, active)
	}
	if pi.BtnPlayer(pi.X, player) {
		pi.Pal(btnX, active)
	}
	if pi.BtnPlayer(pi.Left, player) {
		pi.Pal(left, active)
	}
	if pi.BtnPlayer(pi.Up, player) {
		pi.Pal(up, active)
	}
	if pi.BtnPlayer(pi.Down, player) {
		pi.Pal(down, active)
	}
	if pi.BtnPlayer(pi.Right, player) {
		pi.Pal(right, active)
	}
	// draw gamepad sprite:
	pi.SprSize(0, x, y, 8.0, 4.0)
	// additionally draw a player number (1 light - player 0, 2 lights - player 1)
	drawPlayerNumber(x, y, player)
}

func drawPlayerNumber(x int, y int, player int) {
	pi.PalReset()
	for i := 0; i <= player; i++ {
		pi.Pset(x+50-i*2, y+8, yellow)
	}
}
