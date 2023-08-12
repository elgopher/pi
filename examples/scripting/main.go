// DevTools Terminal lets you write and execute Go code live when your game is running.
// You can run any Pi functions and access any Pi variables. You can also access
// your own variables and functions, but in order to do that you must first export
// them. This example shows how to do that.
//
// Disclaimer: even though you can write and execute Go code when the game is running,
// this does not mean that the Go code is run in parallel. All operations written to
// terminal are run sequentially in the next iteration of the main game loop.
package main

import (
	"embed"
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools"
	"github.com/elgopher/pi/ebitengine"
)

//go:embed sprite-sheet.png
var resources embed.FS

func main() {
	// export variable, so it can be used in terminal, for example:
	//
	// v := variable
	devtools.Export("variable", 1)

	var another = 1

	// export pointer to variable. You can then replace entire variable in terminal:
	//
	// *another = 2
	devtools.Export("another", &another)

	// export function, so it can be used in terminal, for example:
	//
	// drawCircle()
	//
	// Please note, that if you want to see the circle you must first pause the game.
	devtools.Export("drawCircle", drawCircle)

	// export pointer to variable which contains function. You can swap this function with your own code in terminal:
	//
	// *drawCallback = func() { pi.Print("UPDATED", 10, 10, 7) }
	devtools.Export("drawCallback", &drawCallback)

	// export type Struct. It is defined in main therefore you can use something like this:
	//
	// s := Struct{}
	devtools.ExportType[Struct]()

	// export function accepting previously defined struct type:
	//
	// funcAcceptingStruct(Struct{})
	devtools.Export("funcAcceptingStruct", funcAcceptingStruct)

	pi.Load(resources)
	pi.Draw = func() {
		pi.Cls()

		// Animate hello world
		for i := 0; i < 12; i++ {
			x := 20 + i*8
			y := pi.Cos(pi.Time+float64(i)/64) * 60
			pi.Spr(i, x, 60+int(y))
		}

		// run function, which can be replaced in terminal by writing:
		//
		// *drawCallback = func() { pi.Print("UPDATED", 10, 10, 7) }
		drawCallback()

		// you can also replace entire pi.Draw or pi.Update callbacks by writing in terminal:
		//
		// pi.Draw = func() { pi.Cls(); pi.Print("DRAW UPDATED", 10, 10, 7) }
	}

	// Run game with devtools (write pause or p in terminal to pause the game)
	devtools.MustRun(ebitengine.Run)
}

func drawCircle() {
	pi.CircFill(64, 64, 11, 8)
}

var drawCallback = func() {}

type Struct struct {
	Field string
}

func funcAcceptingStruct(p Struct) {
	fmt.Println("Struct received", p)
}
