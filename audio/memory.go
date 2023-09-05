// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

var (
	SFX [64]SoundEffect
	Pat [64]Pattern
)

func Sync() {
	for i, sfx := range SFX {
		system.SetSfx(i, sfx)
	}
	for i, pattern := range Pat {
		system.SetMusic(i, pattern)
	}
}
