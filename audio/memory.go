// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"bytes"
	"fmt"
)

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

// Save stores audio system state to byte slice. State is stored in binary form.
func Save() ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	buffer.WriteByte(schemaVersion)

	for i := 0; i <= maxSfxNo; i++ {
		sfx := Sfx[i]
		for _, note := range sfx.Notes {
			buffer.WriteByte(byte(note.Pitch))
			buffer.WriteByte(byte(note.Instrument))
			buffer.WriteByte(byte(note.Volume))
			buffer.WriteByte(byte(note.Effect))
		}
		buffer.WriteByte(sfx.Speed)
		buffer.WriteByte(sfx.LoopStart)
		buffer.WriteByte(sfx.LoopStop)
		buffer.WriteByte(sfx.Detune)
		buffer.WriteByte(sfx.Reverb)
		buffer.WriteByte(sfx.Dampen)
		buffer.WriteByte(boolToByte(sfx.Noiz))
		buffer.WriteByte(boolToByte(sfx.Buzz))
	}

	for i := 0; i <= maxPatternNo; i++ {
		pattern := Pat[i]
		for _, sfx := range pattern.Sfx {
			buffer.WriteByte(sfx.SfxNo)
			buffer.WriteByte(boolToByte(sfx.Enabled))
		}
		buffer.WriteByte(boolToByte(pattern.BeginLoop))
		buffer.WriteByte(boolToByte(pattern.EndLoop))
		buffer.WriteByte(boolToByte(pattern.StopAtTheEnd))
	}

	return buffer.Bytes(), nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}

	return 0
}

// Load restores audio system state from byte slice. State is restored from binary form.
func Load(state []byte) error {
	if len(state) == 0 {
		return fmt.Errorf("state is empty")
	}

	version := state[0]
	if schemaVersion != version {
		return fmt.Errorf("state version %d is not supported. Only %d is supported.", version, schemaVersion)
	}

	const expectedStateLen = 9409
	if len(state) != expectedStateLen {
		return fmt.Errorf("invalid length of state. Must be %d.", expectedStateLen)
	}

	offset := 1

	for sfxNo := 0; sfxNo <= maxSfxNo; sfxNo++ {
		var sfx SoundEffect

		for j, note := range sfx.Notes {
			note.Pitch = Pitch(state[offset])
			offset++
			note.Instrument = Instrument(state[offset])
			offset++
			note.Volume = Volume(state[offset])
			offset++
			note.Effect = Effect(state[offset])
			offset++

			sfx.Notes[j] = note
		}

		sfx.Speed = state[offset]
		offset++
		sfx.LoopStart = state[offset]
		offset++
		sfx.LoopStop = state[offset]
		offset++
		sfx.Detune = state[offset]
		offset++
		sfx.Reverb = state[offset]
		offset++
		sfx.Dampen = state[offset]
		offset++
		sfx.Noiz = byteToBool(state[offset])
		offset++
		sfx.Buzz = byteToBool(state[offset])
		offset++

		Sfx[sfxNo] = sfx
	}

	for patterNo := 0; patterNo <= maxPatternNo; patterNo++ {
		var pattern Pattern

		for j, sfx := range pattern.Sfx {
			sfx.SfxNo = state[offset]
			offset++
			sfx.Enabled = byteToBool(state[offset])
			offset++

			pattern.Sfx[j] = sfx
		}

		pattern.BeginLoop = byteToBool(state[offset])
		offset++
		pattern.EndLoop = byteToBool(state[offset])
		offset++
		pattern.StopAtTheEnd = byteToBool(state[offset])
		offset++

		Pat[patterNo] = pattern
	}

	return nil
}

func byteToBool(b byte) bool {
	return b == 1
}
