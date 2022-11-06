package pi

func Audio() AudioSystem {
	return audio
}

var audio = AudioSystem{
	Effects: make([]SoundEffect, 64),
}

func init() {
	for i := 0; i < len(audio.Effects); i++ {
		audio.Effects[i].Notes = make([]Note, 32)
	}
}

type AudioSystem struct {
	Effects []SoundEffect // TODO load sound effects from audio.json
	Plan    AudioPlan
}

type SoundEffect struct {
	No    byte
	Speed byte
	Notes []Note
}

type Note struct {
	Pitch  uint16
	Volume uint // 0 - 7
	Wave   Wave
}

type Wave byte

const (
	WaveTriangle  Wave = 1
	WaveTiltedSaw Wave = 2
	WaveSaw       Wave = 2
	WaveSquare    Wave = 3
)

type AudioPlan struct {
	Channel [4]ChannelPlan
}

type ChannelPlan struct {
	Time               float64
	SoundEffect        *SoundEffect
	NoteStart, NoteEnd int
}

// Sfx plays sound effect with given number on any available channel.
func Sfx(no int) {
	e := audio.Effects[no]
	audio.Plan.Channel[anyChannel()] = ChannelPlan{
		Time:        Time(),
		SoundEffect: &e,
	}
}

func anyChannel() int {
	for no, ch := range audio.Plan.Channel {
		if ch.SoundEffect == nil {
			return no
		}
	}

	return 0 // TODO always 0?
}

// SfxChan plays sound effect with given number on specified channel.
func SfxChan(no byte, channel Chan) {
	e := audio.Effects[no]
	audio.Plan.Channel[channel.number()] = ChannelPlan{
		Time:        Time(),
		SoundEffect: &e,
	}
}

type Chan int

const ChanAny Chan = -1

func (c Chan) number() int {
	if c == ChanAny {
		return anyChannel()
	}

	if c < ChanAny || c > 3 {
		return 0
	}

	return int(c)
}

func SfxChanRange(no byte, channel Chan, noteStart, noteEnd int) {
	e := audio.Effects[no]
	audio.Plan.Channel[channel.number()] = ChannelPlan{
		Time:        Time(),
		SoundEffect: &e,
		NoteStart:   noteStart,
		NoteEnd:     noteEnd,
	}
}

func SfxStop(no byte) {
	for _, ch := range audio.Plan.Channel {
		if ch.SoundEffect.No == no {
			ch.SoundEffect = nil
		}
	}
}
