package pi

import "time"

var timeSeconds float64
var timeStarted time.Time

// Time returns the amount of time since game was run, as a (fractional) number of seconds.
//
// Calling Time() multiple times in the same frame will always return the same result.
func Time() float64 {
	return timeSeconds
}

func updateTime() {
	timePassed := time.Since(timeStarted)
	timeSeconds = float64(timePassed) / float64(time.Second)
}
