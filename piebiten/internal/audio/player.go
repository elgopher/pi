// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"log"
	"math"
	"slices"
	"sort"
	"sync"
	"unsafe"

	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/pimath"
)

const chanLen = 4

func newPlayer() *player {
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
		channels: [chanLen]channel{
			defaultChannel, defaultChannel, defaultChannel, defaultChannel,
		},
	}
}

type player struct {
	mutex         sync.Mutex
	samplesByAddr map[uintptr]*piaudio.Sample
	channels      [chanLen]channel

	commandsByTime [chanLen][]command // each channel's planned commands sorted by time

	currentTime float64
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

	for i := 0; i < n; i += 4 {
		p.currentTime += sampleTime
		p.runCommands()
		p.read(out[i : i+4])
	}

	return
}

func (p *player) runCommands() {
	for i := 0; i < chanLen; i++ {
		selectedChan := &p.channels[i]

		processed := 0

		for _, cmd := range p.commandsByTime[i] {
			if cmd.time > p.currentTime {
				break
			}

			switch cmd.kind {
			case cmdKindSetSample:
				switch {
				case cmd.sampleAddr == 0:
					selectedChan.active = false
					selectedChan.sampleData = nil
				case p.samplesByAddr[cmd.sampleAddr] == nil:
					log.Printf("[piaudio] SetSample failed: Sample not found, addr: 0x%x", cmd.sampleAddr)
					selectedChan.active = false
					selectedChan.sampleData = nil
				default:
					selectedChan.active = true
					sample := p.samplesByAddr[cmd.sampleAddr]
					selectedChan.sampleData = sample.Data()
					selectedChan.sampleRate = sample.SampleRate()
				}
				selectedChan.position = float64(cmd.offset)
			case cmdKindSetLoop:
				selectedChan.loop = cmd.loop
			case cmdKindSetPitch:
				selectedChan.pitch = cmd.pitch
			case cmdKindSetVolume:
				selectedChan.volume = cmd.vol
			case cmdKindClearChan:
				// ClearChan was already called in SendCommands
			}
			processed++
		}

		copy(p.commandsByTime[i], p.commandsByTime[i][processed:])
		p.commandsByTime[i] = p.commandsByTime[i][:len(p.commandsByTime[i])-processed]
	}
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

func (p *player) SendCommands(cmds []command) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, cmd := range cmds {
		if cmd.time < p.currentTime {
			log.Printf("Discarding late audio command with time %f, but current time is %f", cmd.time, p.currentTime)
			continue
		}
		if cmd.kind == cmdKindClearChan {
			p.clearChan(cmd.ch, cmd.time)
			continue
		}
		for i := 0; i < chanLen; i++ {
			chanNum := piaudio.Chan(1 << i)
			// a single command can be executed on multiple channels at once
			if cmd.ch&chanNum != 0 {
				p.commandsByTime[i] = append(p.commandsByTime[i], cmd)
			}
		}
	}

	for _, commands := range p.commandsByTime {
		// sort again by time, because new commands may have been inserted between existing ones
		sort.SliceStable(commands, func(i, j int) bool {
			return commands[i].time < commands[j].time
		})
	}
}

func (p *player) clearChan(ch piaudio.Chan, time float64) {
	for i := 0; i < chanLen; i++ {
		chanNum := piaudio.Chan(1 << i)
		if ch&chanNum != 0 {
			idx := sort.Search(len(p.commandsByTime[i]), func(idx int) bool {
				cmd := p.commandsByTime[i][idx]
				return cmd.time >= time
			})
			p.commandsByTime[i] = p.commandsByTime[i][:idx]
		}
	}
}

func (p *player) CurrentTime() float64 {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.currentTime
}

const sampleTime = 1.0 / float64(CtxSampleRate)
