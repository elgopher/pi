// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package input

func IsPressedRepeatably(duration uint) bool {
	const (
		pressDuration = 15
		pressInterval = 4
	)

	if duration == 1 {
		return true
	}

	return duration >= pressDuration+1 && duration%pressInterval == 0
}
