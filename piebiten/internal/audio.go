// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/pimath"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"log"
	"math"
	"slices"
	"sort"
	"sync"
	"time"
	"unsafe"
)

const CtxSampleRate = 44100

func StartAudioBackend(ctx *audio.Context) *AudioBackend {
	timeFromPlayer := make(chan float64, 100)

	thePlayer := newPlayer(timeFromPlayer)
	ebitenPlayer, err := ctx.NewPlayer(thePlayer)
	if err != nil {
		panic("failed to create Ebitengine player: " + err.Error())
	}
	ebitenPlayer.SetBufferSize(time.Duration(audioBufferSizeInSeconds * float64(time.Second)))

	b := &AudioBackend{
		ctx:            ctx,
		timeFromPlayer: timeFromPlayer,
		player:         thePlayer,
		ebitenPlayer:   ebitenPlayer,
	}

	return b
}

// AudioBackend also works in browsers but may glitch if the garbage collector
// blocks the main thread for too long. The backend should ideally use the
// AudioWorklet API to avoid audio glitches.
type AudioBackend struct {
	ctx            *audio.Context
	timeFromPlayer chan float64
	commands       []command
	currentTime    float64
	player         *player
	ebitenPlayer   *audio.Player
}

func (b *AudioBackend) LoadSample(sample *piaudio.Sample) {
	b.player.LoadSample(sample)
}

func (b *AudioBackend) UnloadSample(sample *piaudio.Sample) {
	b.player.UnloadSample(sample)
}

func (b *AudioBackend) scheduleTime(delay float64) float64 {
	return b.currentTime + delay + audioBufferSizeInSeconds
}

func (b *AudioBackend) SetSample(ch piaudio.Chan, sample *piaudio.Sample, offset int, delay float64) {
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

func (b *AudioBackend) SetLoop(ch piaudio.Chan, start, length int, loopType piaudio.LoopType, delay float64) {
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

func (b *AudioBackend) ClearChan(ch piaudio.Chan, delay float64) {
	b.commands = append(b.commands,
		command{
			kind: cmdKindClearChan,
			ch:   ch,
			time: b.scheduleTime(delay),
		},
	)
}

func (b *AudioBackend) SetPitch(ch piaudio.Chan, pitch float64, delay float64) {
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

func (b *AudioBackend) SetVolume(ch piaudio.Chan, vol float64, delay float64) {
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

func (b *AudioBackend) OnBeforeUpdate() {
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

func (b *AudioBackend) OnAfterUpdate() {
	b.player.SendCommands(commandBatch{
		time: b.currentTime,
		cmds: b.commands,
	})
	b.commands = b.commands[:0]
}

func newPlayer(timeFromPlayer chan float64) *player {
	defaultChannel := channel{
		pitch:  1.0,
		volume: 1.0,
		loop: loop{
			stop:     math.MaxInt32,
			loopType: piaudio.LoopNone,
		},
	}
	return &player{
		samplesByAddr: map[uintptr]*piaudio.Sample{},
		time:          timeFromPlayer,
		channels: [4]channel{
			defaultChannel, defaultChannel, defaultChannel, defaultChannel,
		},
	}
}

type player struct {
	mutex         sync.Mutex
	samplesByAddr map[uintptr]*piaudio.Sample
	channels      [4]channel

	currentTime float64

	commandsByTime []command // all planned commands
	time           chan float64
}

type channel struct {
	active     bool
	sampleData []int8
	position   float64 // float for fractional pitch
	pitch      float64
	sampleRate uint16
	volume     float64
	loop       loop
}

func (c *channel) nextSample() (float64, bool) {
	if !c.active || c.sampleData == nil || c.volume <= 0 {
		return 0, false
	}

	pos := int(c.position)
	if pos >= min(len(c.sampleData), c.loop.stop) {
		// End of sample
		if c.loop.loopType == piaudio.LoopForward {
			c.position = float64(c.loop.start)
			pos = c.loop.start
		} else {
			c.active = false
			return 0, false
		}
	}

	sample := float64(c.sampleData[pos])

	// Advance position
	c.position += (float64(c.sampleRate) / CtxSampleRate) * c.pitch

	// Apply volume
	sample *= c.volume

	return sample, true
}

func (p *player) LoadSample(sample *piaudio.Sample) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.samplesByAddr[getPointerAddr(sample)] = piaudio.NewSample(slices.Clone(sample.Data()), sample.SampleRate())
}

func (p *player) UnloadSample(sample *piaudio.Sample) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	delete(p.samplesByAddr, getPointerAddr(sample))
}

func getPointerAddr(sample *piaudio.Sample) uintptr {
	return uintptr(unsafe.Pointer(sample))
}

// Read is called by Ebitengine from a separate goroutine.
// out contains 16-bit stereo PCM data in little-endian format.
func (p *player) Read(out []byte) (n int, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	n = len(out)
	p.time <- p.currentTime + sampleTime*float64(n)

	for i := 0; i < n; i += 4 {
		p.currentTime += sampleTime
		p.runCommands()
		p.read(out[i : i+4])
	}

	return
}

func (p *player) runCommands() {
	processed := 0

	for _, cmd := range p.commandsByTime {
		if cmd.time > p.currentTime {
			break
		}

		c := &p.channels[cmd.ch]

		switch cmd.kind {
		case cmdKindSetSample:
			switch {
			case cmd.sampleAddr == 0:
				c.active = false
				c.sampleData = nil
			case p.samplesByAddr[cmd.sampleAddr] == nil:
				log.Printf("[piaudio] SetSample failed: Sample not found, addr: 0x%x", cmd.sampleAddr)
				c.active = false
				c.sampleData = nil
			default:
				c.active = true
				sample := p.samplesByAddr[cmd.sampleAddr]
				c.sampleData = sample.Data()
				c.sampleRate = sample.SampleRate()
			}
			c.position = float64(cmd.offset)
		case cmdKindSetLoop:
			c.loop = cmd.loop
		case cmdKindSetPitch:
			c.pitch = cmd.pitch
		case cmdKindSetVolume:
			c.volume = cmd.vol
		case cmdKindClearChan:
			// ClearChan was already called in SendCommands
		}
		processed++
	}

	copy(p.commandsByTime, p.commandsByTime[processed:])
	p.commandsByTime = p.commandsByTime[:len(p.commandsByTime)-processed]
}

func (p *player) read(out []byte) {
	numSamples := len(out) / 4

	for i := 0; i < numSamples; i++ {
		var mixL, mixR float64 // -128..127

		for ch := 0; ch < len(p.channels); ch++ {
			sample, ok := p.channels[ch].nextSample()
			if !ok {
				continue
			}

			// Mix
			if ch == 0 || ch == 3 {
				mixL += sample
			} else {
				mixR += sample
			}
		}

		// Write stereo int16 LE PCM
		writeInt16LE(out[i*4:], mixL)
		writeInt16LE(out[i*4+2:], mixR)
	}
}

// val must be [-128..127]
func writeInt16LE(out []byte, val float64) {
	// Scale from [-128..127] to [-32768..32767]
	val16 := int32(val * 256)

	val16 = pimath.Clamp(val16, -32768, 32767)

	sample := int16(val16)

	// Write little-endian
	out[0] = byte(sample)
	out[1] = byte(sample >> 8)
}

type commandBatch struct {
	time float64
	cmds []command
}

func (p *player) SendCommands(batch commandBatch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, cmd := range batch.cmds {
		if cmd.time < p.currentTime {
			log.Printf("Discarding late audio command with time %f, but current time is %f", cmd.time, p.currentTime)
			continue
		}
		if cmd.kind == cmdKindClearChan {
			p.clearChan(cmd.ch, cmd.time)
			continue
		}
		p.commandsByTime = append(p.commandsByTime, cmd)
	}

	// sort again by time, because new commands may have been inserted between existing ones
	sort.SliceStable(p.commandsByTime, func(i, j int) bool {
		return p.commandsByTime[i].time < p.commandsByTime[j].time
	})
}

// clearChan is O(n^2).
// It could be optimized to use a separate command list for each channel.
// Then complexity will be O(n)
func (p *player) clearChan(ch piaudio.Chan, time float64) {
	for j := len(p.commandsByTime) - 1; j >= 0; j-- {
		cmd := p.commandsByTime[j]
		noMoreCommands := cmd.time < time
		if noMoreCommands {
			return
		}
		if cmd.ch == ch {
			// remove cmd
			copy(p.commandsByTime[j:], p.commandsByTime[j+1:])
			p.commandsByTime = p.commandsByTime[:len(p.commandsByTime)-1]
		}
	}
}

const sampleTime = 1.0 / float64(CtxSampleRate)
