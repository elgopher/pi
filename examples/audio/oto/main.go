package main

import (
	"io"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/oto/v2"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
	"github.com/elgopher/pi/key"
)

var (
	sampleRate      = 44100
	channelCount    = 2
	bitDepthInBytes = 2
	freq            = 400
	m               sync.Mutex
)

type SineWave struct {
	length int64
	pos    int64

	remaining []byte
}

func NewSineWave(freq float64, duration time.Duration) *SineWave {
	l := int64(channelCount) * int64(bitDepthInBytes) * int64(sampleRate) * int64(duration) / int64(time.Second)
	l = l / 4 * 4
	return &SineWave{
		length: l,
	}
}

func (s *SineWave) Read(buf []byte) (int, error) {
	m.Lock()
	defer m.Unlock()

	//fmt.Println("READ ", len(buf), time.Now().UnixMilli())
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		copy(s.remaining, s.remaining[n:])
		s.remaining = s.remaining[:len(s.remaining)-n]
		return n, nil
	}

	if s.pos == s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) > s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := float64(sampleRate) / float64(freq)

	num := (bitDepthInBytes) * (channelCount)
	p := s.pos / int64(num)
	switch bitDepthInBytes {
	case 1:
		for i := 0; i < len(buf)/num; i++ {
			const max = 127
			b := int(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < channelCount; ch++ {
				buf[num*i+ch] = byte(b + 128)

			}
			p++
		}
	case 2:
		for i := 0; i < len(buf)/num; i++ {
			const max = 32767
			phase := 2 * math.Pi * float64(p) / length
			fff := math.Sin(phase) * 0.3
			b := int16(fff * max)
			for ch := 0; ch < channelCount; ch++ {
				buf[num*i+2*ch] = byte(b)
				buf[num*i+1+2*ch] = byte(b >> 8)
			}
			p++
		}
	}

	s.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		s.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}

func play(context *oto.Context, freq float64, duration time.Duration) oto.Player {
	p := context.NewPlayer(NewSineWave(freq, duration))
	p.(oto.BufferSizeSetter).SetBufferSize(8192)
	p.SetVolume(0.05)
	p.Play()
	return p
}

func run() (oto.Player, error) {
	const (
		freqC = 523.3
		freqE = 659.3
		freqG = 784.0
		freqA = 850.0
	)

	c, ready, err := oto.NewContext(sampleRate, channelCount, bitDepthInBytes)
	if err != nil {
		return nil, err
	}
	<-ready

	player := play(c, freqG, 20*time.Second)

	return player, nil
}

func main() {
	//flag.Parse()
	player, _ := run()
	pi.Update = func() {
		player.Play()
		if key.Btnp(key.Space) {
			if freq == 400 {
				freq = 0
			} else if freq == 0 {
				freq = 500
			} else if freq == 500 {
				freq = 0
			}
		}

	}
	ebitengine.MustRun()

}
