// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/pimath"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"time"
)

const CtxSampleRate = 44100

func StartAudioBackend(ctx *audio.Context) *Backend {
	timeFromPlayer := make(chan float64, 100)

	thePlayer := newPlayer(timeFromPlayer)
	ebitenPlayer, err := ctx.NewPlayer(thePlayer)
	if err != nil {
		panic("failed to create Ebitengine player: " + err.Error())
	}
	ebitenPlayer.SetBufferSize(time.Duration(audioBufferSizeInSeconds * float64(time.Second)))

	b := &Backend{
		ctx:            ctx,
		timeFromPlayer: timeFromPlayer,
		player:         thePlayer,
		ebitenPlayer:   ebitenPlayer,
	}

	return b
}

// Backend also works in browsers but may glitch if the garbage collector
// blocks the main thread for too long. The backend should ideally use the
// AudioWorklet API to avoid audio glitches.
type Backend struct {
	ctx            *audio.Context
	timeFromPlayer chan float64
	commands       []command
	currentTime    float64
	player         *player
	ebitenPlayer   *audio.Player
}

func (b *Backend) LoadSample(sample *piaudio.Sample) {
	b.player.LoadSample(sample)
}

func (b *Backend) UnloadSample(sample *piaudio.Sample) {
	b.player.UnloadSample(sample)
}

func (b *Backend) scheduleTime(delay float64) float64 {
	return b.currentTime + delay + audioBufferSizeInSeconds
}

func (b *Backend) SetSample(ch piaudio.Chan, sample *piaudio.Sample, offset int, delay float64) {
	b.commands = append(b.commands,
		command{
			kind:       cmdKindSetSample,
			ch:         ch,
			sampleAddr: getPointerAddr(sample),
			offset:     offset,
			time:       b.scheduleTime(delay),
		},
	)
}

type loop struct {
	start, stop int
	loopType    piaudio.LoopType
}

func (b *Backend) SetLoop(ch piaudio.Chan, start, length int, loopType piaudio.LoopType, delay float64) {
	b.commands = append(b.commands,
		command{
			kind: cmdKindSetLoop,
			ch:   ch,
			loop: loop{
				start:    start,
				stop:     start + length - 1,
				loopType: loopType,
			},
			time: b.scheduleTime(delay),
		},
	)
}

func (b *Backend) ClearChan(ch piaudio.Chan, delay float64) {
	b.commands = append(b.commands,
		command{
			kind: cmdKindClearChan,
			ch:   ch,
			time: b.scheduleTime(delay),
		},
	)
}

func (b *Backend) SetPitch(ch piaudio.Chan, pitch float64, delay float64) {
	if pitch < 0 {
		pitch = 0
	}
	b.commands = append(b.commands,
		command{
			kind:  cmdKindSetPitch,
			ch:    ch,
			pitch: pitch,
			time:  b.scheduleTime(delay),
		},
	)
}

func (b *Backend) SetVolume(ch piaudio.Chan, vol float64, delay float64) {
	vol = pimath.Clamp(vol, 0, 1)

	b.commands = append(b.commands,
		command{
			kind: cmdKindSetVolume,
			ch:   ch,
			vol:  vol,
			time: b.scheduleTime(delay),
		},
	)
}

type cmdKind string

const (
	cmdKindSetSample cmdKind = "setSample"
	cmdKindSetLoop   cmdKind = "setLoop"
	cmdKindClearChan cmdKind = "clearChan"
	cmdKindSetPitch  cmdKind = "setPitch"
	cmdKindSetVolume cmdKind = "setVolume"
)

type command struct {
	kind       cmdKind
	ch         piaudio.Chan
	sampleAddr uintptr
	offset     int
	pitch      float64
	time       float64
	vol        float64
	loop       loop
}

func (b *Backend) OnBeforeUpdate() {
	if !b.ebitenPlayer.IsPlaying() {
		b.ebitenPlayer.Play()
	}

	for {
		select {
		case t := <-b.timeFromPlayer:
			b.currentTime = t
			piaudio.Time = t
		default:
			return
		}
	}
}

func (b *Backend) OnAfterUpdate() {
	b.player.SendCommands(b.commands)
	b.commands = b.commands[:0]
}
