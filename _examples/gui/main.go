// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// This example demonstrates building a simple GUI hierarchy with pigui.
// It shows:
//   - A panel (container) with a local coordinate system
//   - Three buttons arranged vertically inside the panel
//   - Clicking a button logs its label
//
// The layout is recalculated relative to the panel's position,
// showing pigui’s tree-structure approach.
package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pigui"
	"log"
)

// colors used in this example (default Pi palette):
const (
	lightBlue = 28
	white     = 7
	darkBlue  = 1
	lightGray = 6
	blue      = 12
)

func main() {
	pi.SetScreenSize(128, 128)
	// create the root of the entire GUI element tree
	root := pigui.New()
	// add a panel (container) at global coordinates
	panel := attachPanel(root, 32, 32, 63, 63)
	// add buttons to the panel using its local coordinate system
	attachButton(panel, 10, 9, 44, 14, "BUTTON 1")
	attachButton(panel, 10, 25, 44, 14, "BUTTON 2")
	// add a button with a callback that runs when the user clicks
	// and releases the left mouse button while staying inside its area
	btn3 := attachButton(panel, 10, 41, 44, 14, "BUTTON 3")
	btn3.OnTap = func(event pigui.Event) {
		log.Println("Button 3 was tapped")
	}

	pi.Update = func() {
		// root.Update() must be called in the game loop
		root.Update()
	}

	pi.Draw = func() {
		pi.Cls()
		// root.Draw() must be called in the game loop
		root.Draw()
	}

	piebiten.Run()
}

func attachPanel(parent *pigui.Element, x, y, w, h int) *pigui.Element {
	panel := pigui.Attach(parent, x, y, w, h)
	panel.OnDraw = func(event pigui.DrawEvent) {
		pi.SetColor(lightBlue)
		pi.Rect(0, 0, panel.W-1, panel.H-1)
		pi.SetColor(darkBlue)
		pi.RectFill(1, 1, panel.W-2, panel.H-2)
	}
	return panel
}

func attachButton(parent *pigui.Element, x, y, w, h int, label string) *pigui.Element {
	btn := pigui.Attach(parent, x, y, w, h)
	btn.OnDraw = func(event pigui.DrawEvent) {
		var frame, bg, text pi.Color = lightGray, blue, white
		if event.HasPointer {
			frame, bg, text = lightGray, lightBlue, white
		}

		if event.Pressed {
			pi.Camera.Y -= 1 // the camera is automatically reset after drawing the element
			bg = blue
		}

		pi.SetColor(frame)
		pi.Rect(0, 0, w-2, h-2)

		pi.SetColor(bg)
		pi.RectFill(1, 1, w-3, h-3)

		pi.SetColor(text)
		picofont.Print(label, 6, 4)
	}
	return btn
}
