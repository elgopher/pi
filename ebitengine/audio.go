package ebitengine

import (
	"fmt"

	"github.com/hajimehoshi/oto/v2"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio"
)

type otoAudio struct {
	context *oto.Context
	player  oto.Player
	buffer  *audio.ReaderBuffer
}

func (a *otoAudio) Start() error {
	if a.context == nil {
		audioCtx, ready, err := oto.NewContext(audio.SampleRate, 1, audio.BitDepth/8)
		if err != nil {
			return fmt.Errorf("problem creating Oto Audio Context: %w", err)
		}
		<-ready

		a.context = audioCtx

		a.buffer = &audio.ReaderBuffer{}
		player := audioCtx.NewPlayer(a.buffer)
		player.(oto.BufferSizeSetter).SetBufferSize(2940)
		a.player = player
		go a.player.Play()
	}

	return nil
}

func (a *otoAudio) Update() {
	audio.Write(pi.Time(), pi.Audio().Plan, a.buffer)
}

func (a *otoAudio) Stop() {
	if a.context != nil {
		a.player.Close()
		a.context.Suspend()
	}
}
