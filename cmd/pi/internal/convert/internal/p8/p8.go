package p8

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/icza/bitio"

	"github.com/elgopher/pi/audio"
)

func ConvertToAudioSfx(inputFile, outputFile string) error {
	parser := Parser{}
	file, err := parser.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("error parsing p8 file %s: %w", inputFile, err)
	}

	for _, section := range file.Sections {
		if section.Name == "__sfx__" {
			sfx, err := decodeSfx(section.Lines)
			if err != nil {
				return fmt.Errorf("error decoding __sfx__ section from p8 file %s: %w", inputFile, err)
			}
			audio.Sfx = sfx
		}
	}

	bytes, err := audio.Save()
	if err != nil {
		return fmt.Errorf("saving audio failed: %w", err)
	}

	err = os.WriteFile(outputFile, bytes, 0644)
	if err != nil {
		return fmt.Errorf("writing %s failed: %w", outputFile, err)
	}

	return nil
}

func decodeSfx(lines []string) (sfx [64]audio.SoundEffect, err error) {
	// each line is a sound effect
	for no, line := range lines {
		decoded, err := hex.DecodeString(line)
		if err != nil {
			return sfx, err
		}

		notes, err := decodeSfxNotes(decoded[4:])
		if err != nil {
			return sfx, err
		}

		editorModeAndFilters := bitio.NewReader(bytes.NewReader(decoded[0:1]))
		editorMode, err := editorModeAndFilters.ReadBool()
		if err != nil {
			return sfx, err
		}
		_ = editorMode

		noiz, err := editorModeAndFilters.ReadBool()
		if err != nil {
			return sfx, err
		}

		buzz, err := editorModeAndFilters.ReadBool()
		if err != nil {
			return sfx, err
		}

		detune, err := editorModeAndFilters.ReadBits(2)
		if err != nil {
			return sfx, err
		}

		reverb, err := editorModeAndFilters.ReadBits(2)
		if err != nil {
			return sfx, err
		}

		dampen, err := editorModeAndFilters.ReadBits(2)
		if err != nil {
			return sfx, err
		}

		sfx[no] = audio.SoundEffect{
			Speed:     decoded[1],
			LoopStart: decoded[2],
			LoopStop:  decoded[3],
			Notes:     notes,
			Noiz:      noiz,
			Buzz:      buzz,
			Detune:    byte(detune),
			Reverb:    byte(reverb),
			Dampen:    byte(dampen),
		}
	}

	return
}

func decodeSfxNotes(b []byte) (notes [32]audio.Note, err error) {
	reader := bitio.NewReader(bytes.NewBuffer(b))
	for i := 0; i < 32; i++ {
		pitch, err := reader.ReadByte()
		if err != nil {
			return notes, err
		}
		notes[i].Pitch = audio.Pitch(pitch)

		waveform, err := reader.ReadBits(4)
		if err != nil {
			return notes, err
		}
		notes[i].Instrument = audio.Instrument(waveform)

		volume, err := reader.ReadBits(4)
		if err != nil {
			return notes, err
		}
		notes[i].Volume = audio.Volume(volume)

		effect, err := reader.ReadBits(4)
		if err != nil {
			return notes, err
		}
		notes[i].Effect = audio.Effect(effect)
	}

	return
}
