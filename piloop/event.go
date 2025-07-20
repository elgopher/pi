// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piloop

type Event string

const (
	EventInit        Event = "init"         // when the game is started, just before the first frame
	EventFrameStart  Event = "frame_start"  // beginning of the frame
	EventUpdate      Event = "update"       // after pi.Update
	EventLateUpdate  Event = "late_update"  // after EventUpdate
	EventDraw        Event = "draw"         // after pi.Draw
	EventLateDraw    Event = "late_draw"    // after EventDraw
	EventWindowClose Event = "window_close" // when a user closes the window (desktop only)
)
