// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piloop

type Event string

const (
	EventGameStarted  Event = "game_started"  // just before first frame
	EventFrameStarted Event = "frame_started" // beginning of the frame
	EventUpdate       Event = "update"        // after pi.Update
	EventLateUpdate   Event = "late_update"   // after EventUpdate
	EventDraw         Event = "draw"          // after pi.Draw
	EventLateDraw     Event = "late_draw"     // after EventDraw
	EventWindowClosed Event = "window_closed" // when user closes the window (desktop only)
)
