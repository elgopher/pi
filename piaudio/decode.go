// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"
)

// DecodeRaw decodes uncompressed raw data (no headers) into a *Sample.
// Expects 8-bit mono PCM with samples as int8 (-128..127).
func DecodeRaw(raw []byte) *Sample {
	data := byteSliceToInt8Slice(raw)

	return NewSample(data)
}

func byteSliceToInt8Slice(b []byte) []int8 {
	return unsafe.Slice((*int8)(unsafe.Pointer(unsafe.SliceData(b))), len(b))
}

// DecodeWav decodes a WAV file into a *Sample, panicking if decoding fails.
// Expects 8-bit mono PCM with samples as int8 (-128..127).
func DecodeWav(wav []byte) *Sample {
	sample, err := DecodeWavOrErr(wav)
	if err != nil {
		panic(err)
	}
	return sample
}

// DecodeWavOrErr decodes a WAV file into a *Sample or returns an error if decoding fails.
// Expects 8-bit mono PCM with samples as int8 (-128..127).
func DecodeWavOrErr(wav []byte) (*Sample, error) {
	if len(wav) < 44 {
		return nil, errors.New("WAV too short")
	}

	if string(wav[0:4]) != "RIFF" || string(wav[8:12]) != "WAVE" {
		return nil, errors.New("not a valid WAV file")
	}

	var bitsPerSample uint16
	var numChannels uint16
	var dataChunk []byte

	// WAV chunk parsing
	offset := 12 // skip RIFF header

	for offset+8 <= len(wav) {
		chunkID := string(wav[offset : offset+4])
		chunkSize := int(binary.LittleEndian.Uint32(wav[offset+4 : offset+8]))
		chunkDataStart := offset + 8
		chunkDataEnd := chunkDataStart + chunkSize

		if chunkDataEnd > len(wav) {
			return nil, errors.New("chunk size out of bounds")
		}

		switch chunkID {
		case "fmt ":
			if chunkSize < 16 {
				return nil, errors.New("fmt chunk too small")
			}
			audioFormat := binary.LittleEndian.Uint16(wav[chunkDataStart : chunkDataStart+2])
			if audioFormat != 1 {
				return nil, errors.New("only PCM supported")
			}
			numChannels = binary.LittleEndian.Uint16(wav[chunkDataStart+2 : chunkDataStart+4])
			if numChannels != 1 {
				return nil, fmt.Errorf("only mono supported, got %d channels", numChannels)
			}
			bitsPerSample = binary.LittleEndian.Uint16(wav[chunkDataStart+14 : chunkDataStart+16])

		case "data":
			dataChunk = wav[chunkDataStart:chunkDataEnd]

		default:
			// skip unknown chunk
		}

		offset = chunkDataEnd
		if chunkSize%2 == 1 {
			offset++
		}
	}

	if bitsPerSample == 0 {
		return nil, errors.New("missing fmt chunk")
	}
	if bitsPerSample != 8 {
		return nil, fmt.Errorf("only 8-bit PCM supported, got %d bits", bitsPerSample)
	}
	if len(dataChunk) == 0 {
		return nil, errors.New("no data chunk found")
	}

	// Convert unsigned -> signed
	pcm := make([]int8, len(dataChunk))
	for i, v := range dataChunk {
		pcm[i] = int8(int(v) - 128)
	}

	return NewSample(pcm), nil
}
