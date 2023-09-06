// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

var (
	Sfx [64]SoundEffect // Sound effects
	Pat [64]Pattern     // Music patterns
)

// Sync is required for changes made to Sfx and Pat to be audible.
// Sync is automatically run after each command issued via devtools terminal.
func Sync() {
	for i, sfx := range Sfx {
		system.SetSfx(i, sfx)
	}
	for i, pattern := range Pat {
		system.SetMusic(i, pattern)
	}
}
