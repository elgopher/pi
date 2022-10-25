// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package mem is low-level package for directly accessing memory,
// such as screen pixels, sprite-sheet, fonts or buttons state.
// This data can be manipulated by backend, devtools, utility functions, or
// a game itself. It is very useful for writing custom tools, new backends or
// even entire new API to be used by games. Code using mem package directly
// could be very fast, because it can use low-level Go functions such as copy.
//
// Please note though, that with great power comes great responsibility. You
// can easily shoot yourself in the foot if you are not careful enough how
// you change the data. For example, increasing the SpriteSheetWidth
// without adjusting the SpriteSheetData will likely result in a panic
// during sprite-drawing operations.
package mem

var (
	// TimeSeconds is the number of seconds since game was started
	TimeSeconds float64

	GameLoopStopped bool
)
