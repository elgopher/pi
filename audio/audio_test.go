package audio_test

import (
	"fmt"
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio"
)

func TestWrite(t *testing.T) {
	plan := pi.AudioPlan{}
	buffer := &audio.ReaderBuffer{}

	audio.Write(0.1, plan, buffer)
	fmt.Println(buffer)
}
