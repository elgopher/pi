// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"fmt"
	"sync"
	"time"

	ebitenaudio "github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio"
)

const (
	audioSampleRate = 44100
	channelCount    = 2 // stereo
)

// AudioStream is an abstraction used by ebitengine back-end to consume audio stream generated by the game. The stream
// will be played back to the user. By default, audio package is used. Game developers could use a different audio
// system though. In such case they could set the AudioStream variable to their own implementation.
var AudioStream interface {
	// Read reads generated audio into p buffer.
	//
	// The format is:
	//	[data]      = [sample 0] [sample 1] [sample 2] ...
	//	[sample *]  = [channel left] [channel right] ...
	//	[channel *] = [bits 0-15]
	//
	// Sample rate must be 44100, channel count 2 (stereo) and bit depth 16.
	//
	// Byte ordering is little endian.
	//
	// See [io.Reader] for documentation how to implement this method.
	// When error is returned then game stops with the error.
	Read(p []byte) (n int, err error)
}

func startAudio() (stop func(), ready <-chan struct{}, _ error) {
	if AudioStream == nil {
		// In the web back-end, Audio Worklets will be used. In the beginning, state will be stored to binary form
		// and sent over the MessageChannel to processor. Each call to System methods will send events to processor
		// (again via MessageChannel) instead of directly calling synthesizer. Audio Worklet processor will use
		// pi.Synthesizer. Based on incoming events it will update the pi.Synthesizer.

		state, err := audio.SaveAudio()
		if err != nil {
			return stop, nil, fmt.Errorf("problem saving audio state: %w", err)
		}
		synth := &audio.Synthesizer{}
		if err = synth.Load(state); err != nil {
			return stop, nil, fmt.Errorf("problem loading audio state: %w", err)
		}

		audioSystem := &ebitenPlayerSource{audioSystem: synth}
		audio.SetSystem(audioSystem) // make audio system concurrency-safe
		AudioStream = audioSystem
	}

	audioCtx := ebitenaudio.NewContext(audioSampleRate)

	player, err := audioCtx.NewPlayer(AudioStream)
	if err != nil {
		return func() {}, nil, err
	}
	player.SetBufferSize(60 * time.Millisecond)

	readyChan := make(chan struct{})

	go func() {
		for {
			if audioCtx.IsReady() {
				close(readyChan)

				player.Play()

				return
			}
			time.Sleep(time.Millisecond)
		}
	}()

	return func() {
		_ = player.Close()
	}, readyChan, nil
}

// ebitenPlayerSource implements Ebitengine Player source.
//
// ebitenPlayerSource adds synchronization to all methods.
// Therefore, it can be called concurrently by Ebitegine and the game loop.
type ebitenPlayerSource struct {
	mutex       sync.Mutex
	audioSystem interface {
		audio.System
		ReadSamples(p []float64)
	}

	singleSample   []byte    // singleSample in Ebitengine format - first two bytes left channel, next two bytes right
	remainingBytes int       // number of bytes from singleSample still not copied to p
	floatBuffer    []float64 // reused buffer to avoid allocation on each Read request
}

// reads floats from AudioStream and convert them to Ebitengine format -
// linear PCM (signed 16bits little endian, 2 channel stereo).
func (e *ebitenPlayerSource) Read(p []byte) (int, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	const (
		uint16Bytes = 2
		sampleLen   = channelCount * uint16Bytes
	)

	if len(p) == 0 {
		return 0, nil
	}

	if e.remainingBytes > 0 {
		n := copy(p, e.singleSample[sampleLen-e.remainingBytes:])
		e.remainingBytes = 0
		return n, nil
	}

	if e.singleSample == nil {
		e.singleSample = make([]byte, sampleLen)
	}

	samples := len(p) / sampleLen
	if len(p)%sampleLen != 0 {
		samples += 1
		e.remainingBytes = sampleLen - len(p)%sampleLen
	}

	e.ensureFloatBufferIsBigEnough(samples)

	bytesRead := 0

	e.audioSystem.ReadSamples(e.floatBuffer[:samples])
	for i := 0; i < samples; i++ {
		floatSample := pi.Mid(e.floatBuffer[i], -1, 1)
		sample := int16(floatSample * 0x7FFF) // actually the full int16 range is -0x8000 to 0x7FFF (therefore -0x8000 will never be returned)

		e.singleSample[0] = byte(sample)
		e.singleSample[1] = byte(sample >> 8)
		copy(e.singleSample[2:], e.singleSample[:2]) // copy left to right channel

		copiedBytes := copy(p, e.singleSample)
		p = p[copiedBytes:]
		bytesRead += copiedBytes
	}

	return bytesRead, nil
}

func (e *ebitenPlayerSource) ensureFloatBufferIsBigEnough(size int) {
	if size > len(e.floatBuffer) {
		e.floatBuffer = make([]float64, size)
	}
}

func (e *ebitenPlayerSource) Sfx(sfxNo int, channel audio.Channel, offset, length int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.audioSystem.Sfx(sfxNo, channel, offset, length)
}

func (e *ebitenPlayerSource) Music(patterNo int, fadeMs int, channelMask byte) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.audioSystem.Music(patterNo, fadeMs, channelMask)
}

func (e *ebitenPlayerSource) Stat() audio.Stat {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.audioSystem.Stat()
}

func (e *ebitenPlayerSource) SetSfx(sfxNo int, effect audio.SoundEffect) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.audioSystem.SetSfx(sfxNo, effect)
}

func (e *ebitenPlayerSource) GetSfx(sfxNo int) audio.SoundEffect {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.audioSystem.GetSfx(sfxNo)
}

func (e *ebitenPlayerSource) SetMusic(patterNo int, pattern audio.Pattern) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.audioSystem.SetMusic(patterNo, pattern)
}

func (e *ebitenPlayerSource) GetMusic(patterNo int) audio.Pattern {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.audioSystem.GetMusic(patterNo)
}

func (e *ebitenPlayerSource) Save() ([]byte, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.audioSystem.Save()
}

func (e *ebitenPlayerSource) Load(bytes []byte) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.audioSystem.Load(bytes)
}
