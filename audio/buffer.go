package audio

// Buffer is an overriding circular buffer which stores audio samples.
type Buffer struct {
	samples []float64

	capacity int
	head     int
	tail     int
}

func NewBuffer(length uint) *Buffer {
	return &Buffer{
		samples:  make([]float64, length),
		capacity: 0,
		head:     0,
		tail:     0,
	}
}

func (b *Buffer) Write(in []float64) {
	if len(b.samples) == 0 {
		return
	}

	if len(in) > len(b.samples) {
		in = in[len(in)-len(b.samples):]
	}

	copy(b.samples[b.tail:], in)
	b.tail += len(in)
}

func (b *Buffer) Read(out []float64) (n int) {
	if len(b.samples) == 0 {
		return
	}

	n = b.tail - b.head
	copy(out, b.samples)

	return
}
