//go:build !js

package ebitengine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/elgopher/pi"
)

func startAudio() (stop func(), _ error) {
	if AudioStream == nil {
		AudioStream = pi.Audio()
	}

	audioCtx := audio.NewContext(audioSampleRate)
	player, err := audioCtx.NewPlayer(&audioStreamReader{})
	if err != nil {
		return func() {}, err
	}
	player.SetBufferSize(50 * time.Millisecond)
	player.Play()

	return func() {
		_ = player.Close()
	}, nil
}
