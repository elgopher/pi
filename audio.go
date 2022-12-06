// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Audio returns AudioSystem object.
func Audio() *AudioSystem {
	return audio
}

var audio = &AudioSystem{}

// AudioSystem is safe to use in a concurrent manner.
type AudioSystem struct{}

// Read method is used by back-end to read generated audio stream and play it back to the user. The sample rate is 44100,
// 16 bit depth and stereo (2 audio channels).
//
// Read is (usually) executed concurrently with main game loop. Back-end could decide about buffer size, although
// the higher the size the higher the lag. Usually the buffer is 8KB, which is 46ms of audio.
func (s *AudioSystem) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	// generate silence for now
	for i := 0; i < len(p); i++ {
		p[i] = 0
	}

	return len(p), nil
}
