// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio_test

import (
	_ "embed"
	"testing"

	"github.com/elgopher/pi/audio"
)

func BenchmarkLoad(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := audio.Load(validSave) // 11us
		if err != nil {
			b.Logf("error returned from Synthesizer.Load: %s", err)
			b.Fail()
		}
	}
}

func BenchmarkSynthesizer_Save(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var (
		bytes []byte
		err   error
	)

	for i := 0; i < b.N; i++ {
		bytes, err = audio.Save() // 27us
		if err != nil {
			b.Logf("error returned from Synthesizer.Save: %s", err)
			b.Fail()
		}
	}

	_ = bytes
}
