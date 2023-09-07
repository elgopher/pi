// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio_test

import (
	_ "embed"
	"github.com/elgopher/pi/audio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed "internal/valid-save"
var validSave []byte

func TestSave(t *testing.T) {
	t.Run("should save in binary format", func(t *testing.T) {
		audio.Sfx[3] = validEffect
		audio.Pat[4] = validPattern
		bytes, err := audio.Save()
		require.NoError(t, err)
		assert.Equal(t, validSave, bytes)
	})
}

func TestLoad(t *testing.T) {
	t.Run("should load state", func(t *testing.T) {
		err := audio.Load(validSave)
		require.NoError(t, err)
		assert.Equal(t, validEffect, audio.Sfx[3])
		assert.Equal(t, validPattern, audio.Pat[4])
	})

	t.Run("should return error when state is empty", func(t *testing.T) {
		err := audio.Load([]byte{})
		assert.Error(t, err)
	})

	t.Run("should return error when version is not supported", func(t *testing.T) {
		err := audio.Load([]byte{2})
		assert.Error(t, err)
	})

	t.Run("should return error when state has invalid length", func(t *testing.T) {
		err := audio.Load([]byte{1, 0, 0, 0})
		assert.Error(t, err)
	})
}
