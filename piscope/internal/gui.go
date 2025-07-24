// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pidebug"
	"github.com/elgopher/pi/pigui"
)

func attachToolbar(parent *pigui.Element) *pigui.Element {
	toolbar := pigui.Attach(parent, 0, pi.Screen().H()-9, pi.Screen().W(), 9)
	toolbar.OnDraw = func(event pigui.DrawEvent) {
		prev := pi.SetColor(*bgColor)
		defer pi.SetColor(prev)
		pi.RectFill(0, 0, toolbar.W, toolbar.H)
	}

	// attachIconButton(toolbar, icons.AlignTop, 0) // icon hidden until implemented
	// attachIconButton(toolbar, icons.Screen, 8) // icon hidden for now because screen inspector is the only tab
	// attachIconButton(toolbar, icons.Palette, 16) // icon hidden until implemented
	// attachIconButton(toolbar, icons.Variables, 24) // icon hidden until implemented
	// attachIconButton(toolbar, icons.Paint, 32) // icon hidden until implemented

	snap := attachIconButton(toolbar, icons.Snap, pi.Screen().W()-34)
	snap.OnTap = func(event pigui.Event) {
		captureSnapshot()
	}

	prev := attachIconButton(toolbar, icons.Prev, pi.Screen().W()-24)
	prev.OnTap = func(event pigui.Event) {
		showPrevSnapshot()
	}
	prev.OnUpdate = func(pigui.UpdateEvent) {
		if theScreenRecorder.HasPrev() {
			prev.Icon = icons.Prev
		} else {
			prev.Icon = pi.Sprite{}
		}
	}

	playPause := attachIconButton(toolbar, icons.Pause, pi.Screen().W()-19)
	playPause.OnTap = func(pigui.Event) {
		pauseOrResume()
	}
	playPause.OnUpdate = func(pigui.UpdateEvent) {
		if pidebug.Paused() {
			playPause.Icon = icons.Pause
		} else {
			playPause.Icon = icons.Play
		}
	}

	next := attachIconButton(toolbar, icons.Next, pi.Screen().W()-14)
	next.OnTap = func(event pigui.Event) {
		showNextSnapshot()
	}

	exit := attachIconButton(toolbar, icons.Exit, pi.Screen().W()-8)
	exit.OnTap = func(event pigui.Event) {
		exitConsoleMode()
	}

	return toolbar
}

type IconButton struct {
	*pigui.Element

	Icon pi.Sprite
}

func attachIconButton(parent *pigui.Element, icon pi.Sprite, x int) *IconButton {
	btn := pigui.Attach(parent, x, 0, icon.W, icon.H+1)
	iconBtn := &IconButton{Icon: icon, Element: btn}
	btn.OnDraw = func(event pigui.DrawEvent) {
		y := 0
		if event.Pressed {
			y = 1
		}

		prevColorTable := pi.ColorTables[0]
		defer func() {
			pi.ColorTables[0] = prevColorTable
		}()
		pi.Pal(0, *bgColor) // 0 is bg color in icons.png
		pi.Pal(1, *fgColor) // 1 is fg color in icons.png
		pi.DrawSprite(iconBtn.Icon, 0, y)
	}
	return iconBtn
}
