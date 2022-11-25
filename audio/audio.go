package audio

import (
	"math"
	"sync"
)

const (
	SampleRate = 44100
	Size       = 8192 // ~93ms
	BitDepth   = 16
)

var data = make([]byte, Size)

// Write generates the audio wave based on audioPlan for the next ~93ms, starting at time.
//
// The wave will be written to the buffer at the sample position calculated from time.
// The sample rate of written data is 44100. Data is 16-bit PCM, one channel.
// Byte ordering is little endian.
//
// The data format is as follows:
//
//	[data]      = [sample 0] [sample 1] [sample 2] ...
//	[sample *]  = [byte 0] [byte 1]
func Write(time float64, buffer Buffer) {
	pos := SampleRate * time

	length := float64(SampleRate) / float64(400)

	for i := 0; i < len(data)/2; i++ {
		const max = 32767
		b := int16(math.Sin(2*math.Pi*pos/length) * 0.3 * max)
		data[2*i] = byte(b)
		data[2*i+1] = byte(b >> 8)
		pos++
	}

	buffer.Write(uint64(pos), data)
}

type Buffer interface {
	// Write updates the buffer starting at sample.
	Write(sample uint64, data []byte)
}

// ReaderBuffer is a Buffer which can be also read using io.Reader
//
// Data already read are discarded.
type ReaderBuffer struct {
	start uint64
	data  []byte
	mutex sync.Mutex
}

func (c *ReaderBuffer) Write(sample uint64, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// ensure there is a space
	if len(c.data)+int(c.start) < len(data)+int(sample) {
		enlargedData := make([]byte, len(data)+int(sample)-int(c.start))
		copy(enlargedData, c.data)
		c.data = enlargedData
	}
	if sample < c.start {
		data = data[c.start-sample:]
		sample = c.start
	}
	copy(c.data[sample-c.start:], data)
}

func (c *ReaderBuffer) Read(p []byte) (n int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	n = copy(p, c.data)
	// move data
	copied := copy(c.data, c.data[n:])
	c.data = c.data[:copied]
	c.start += uint64(n)

	return n, nil
}
