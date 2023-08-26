// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio_test

import (
	_ "embed"
	"testing"

	"github.com/elgopher/pi/audio"
)

func BenchmarkSynthesizer_Load(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	s := audio.Synthesizer{}

	for i := 0; i < b.N; i++ {
		err := s.Load(validSave) // 11us
		if err != nil {
			b.Logf("error returned from Synthesizer.Load: %s", err)
			b.Fail()
		}
	}
}

func BenchmarkSynthesizer_Save(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	s := audio.Synthesizer{}

	var (
		bytes []byte
		err   error
	)

	for i := 0; i < b.N; i++ {
		bytes, err = s.Save() // 27us
		if err != nil {
			b.Logf("error returned from Synthesizer.Save: %s", err)
			b.Fail()
		}
	}

	_ = bytes
}
