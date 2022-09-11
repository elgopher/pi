// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package input

func IsPressedRepeatably(duration int) bool {
	const (
		pressDuration = 15 // make it configurable
		pressInterval = 4  // make it configurable
	)

	if duration == 1 {
		return true
	}

	return duration >= pressDuration+1 && duration%pressInterval == 0
}
