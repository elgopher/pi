package main

import (
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pikey"
	"github.com/elgopher/pi/piloop"
)

func main() {
	// stops the game loop when user pressed Escape key
	pikey.RegisterShortcut(piloop.Stop, pikey.Esc)

	piebiten.Run() // this function ends without panic once game loop is stopped

	// Code in here is executed after loop is stopped
}
